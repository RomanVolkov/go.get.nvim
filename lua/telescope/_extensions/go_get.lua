local has_telescope, telescope = pcall(require, "telescope")
local init = require("go_get.init")

if not has_telescope then
	error("This plugins requires nvim-telescope/telescope.nvim")
end

return telescope.register_extension({
	exports = { packages_search = init.packages_search },
})
