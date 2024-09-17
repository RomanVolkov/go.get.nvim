local pickers = require("telescope.pickers")
local previewers = require("telescope.previewers.previewer")
local finders = require("telescope.finders")
local conf = require("telescope.config").values
local actions = require("telescope.actions")
local action_state = require("telescope.actions.state")

local M = {}

function split(str, sep)
	local result = {}
	for match in (str .. sep):gmatch("(.-)" .. sep) do
		table.insert(result, match)
	end
	return result
end

local function loadIndex()
	local fullPath = debug.getinfo(1).source:sub(2)
	local folderPath = fullPath:match("(.*/)")
	if folderPath == nil then
		folderPath = ""
	end
	local file = io.open(folderPath .. "index.txt", "r")

	if not file then
		error("File not found: ")
	end

	local skipLines = 1
	local lines = {}
	local count = 0
	for line in file:lines() do
		count = count + 1
		if count > skipLines then
			-- format will be the same:
			-- url, license, homepage, description
			table.insert(lines, split(line, ";"))
		end
	end

	file:close()
	return lines
end

local installation_queue = {}
process_queue = function()
	local url_value = installation_queue[1]
	if url_value == nil then
		return
	end

	local cmd = "go"
	local cmd_args = { "get", "-u", url_value }
	-- the problem that I'm running it from lua file, not golang project
	-- let's test that
	vim.notify("installing..." .. url_value)

	local uv = vim.loop
	local handle
	local on_exit = function(status)
		if status == 0 then
			vim.notify("installed: " .. url_value)
		else
			vim.notify("failed to install: " .. url_value)
			-- TODO: where to write error?
		end

		uv.close(handle)
		table.remove(installation_queue, 1)
		process_queue()
	end
	handle = uv.spawn(cmd, { args = cmd_args }, on_exit)
end

M.packages_search = function(opts)
	opts = opts or require("telescope.themes").get_dropdown({})
	pickers
		.new(opts, {
			prompt_title = "go get",
			finder = finders.new_table({
				results = loadIndex(),
				entry_maker = function(entry)
					local url = entry[1]
					local license = entry[3]
					if license == nil then
						license = ""
					end
					local description = entry[3]
					if description == nil then
						description = ""
					end
					return {
						value = entry,
						display = entry[1], -- value to display
						-- I don't think that I need a license as search part
						ordinal = entry[1] .. license .. " " .. description, -- value for search
					}
				end,
			}),
			sorter = conf.generic_sorter(opts),
			attach_mappings = function(prompt_bufnr, map)
				actions.select_default:replace(function()
					-- could be a option to close or not automatically
					actions.close(prompt_bufnr)
					local selection = action_state.get_selected_entry()
					-- take the value to use
					-- local url_value = selection["value"][2]
					local url_value = selection["value"][1]
					vim.notify(url_value)
					local current_count = #installation_queue
					installation_queue[#installation_queue + 1] = url_value
					if current_count == 0 then
						process_queue()
					end
				end)
				return true
			end,
			previewer = previewers:new({
				preview_fn = function(self, entry, status)
					local selection = action_state.get_selected_entry()
					local selectionValue = selection["value"]
					-- local url_value = selection["value"][1]
					local bufnr = status.preview_bufnr
					local url = selectionValue[1]
					local license = selectionValue[2]
					local homepage = selectionValue[3]
					local description = selectionValue[4]
					local previewMessage = {
						"Package URL:",
						url,
					}
					if string.len(license) > 0 then
						table.insert(previewMessage, "License:")
						table.insert(previewMessage, license)
					end
					if string.len(homepage) > 0 then
						table.insert(previewMessage, "Homepage")
						table.insert(previewMessage, homepage)
					end
					if string.len(description) > 0 then
						table.insert(previewMessage, "Description")
						table.insert(previewMessage, description)
					end
					vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, previewMessage)
				end,
			}),
		})
		:find()
end

-- M.packages_search()

return M
