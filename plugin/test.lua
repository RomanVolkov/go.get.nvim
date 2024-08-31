local pickers = require("telescope.pickers")
local finders = require("telescope.finders")
local conf = require("telescope.config").values
local actions = require("telescope.actions")
local action_state = require("telescope.actions.state")

local function loadIndex()
	local file = io.open("index.txt", "r")
	if not file then
		error("File not found: ")
	end

	local skipLines = 1
	local lines = {}
	local count = 0
	for line in file:lines() do
		count = count + 1
		if count > skipLines then
			table.insert(lines, { line })
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

local packages_search = function(opts)
	opts = opts or {}
	pickers
		.new(opts, {
			prompt_title = "go get",
			finder = finders.new_table({
				results = loadIndex(),
				entry_maker = function(entry)
					return {
						value = entry,
						display = entry[1], -- value to display
						ordinal = entry[1], -- value for search
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
					local current_count = #installation_queue
					installation_queue[#installation_queue + 1] = url_value
					if current_count == 0 then
						process_queue()
					end
				end)
				return true
			end,
		})
		:find()
end

-- to execute the function
packages_search(require("telescope.themes").get_dropdown({}))

--
