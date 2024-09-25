# `go get` plugin for NeoVim

An extension for [telescope.nvim](https://github.com/nvim-telescope/telescope.nvim) that allows you to install Golang dependencies without leaving of the editor.

### Demo

https://github.com/user-attachments/assets/b0c1ed68-ea01-4455-bdfb-643a22f66a79



## Requirements

- [telescope.nvim](https://github.com/nvim-telescope/telescope.nvim) 


## Setup
You can setup the extension by adding the following to your config:
```lua
require'telescope'.load_extension('go_get')
```

Currently plugin integration is tested with only `lazy` plugin manager, others to be done.

Example of lazy.nvim integration

```lua
return {
  "https://github.com/RomanVolkov/go.get.nvim",
  dependencies = {
    "nvim-telescope/telescope.nvim",
  },
  keys = function()
    return {
      {
        "<leader>gog",
        function()
          require("telescope").extensions.go_get.packages_search()
        end,
        desc = "[G]o [G]et packages",
      },
    }
  end,
}
```

### Capabilities

Currently there is only one function for this project - `packages_search`. It will open Telescope with dropdown menu where you can select a Golang dependency. 
By selecting it (with Enter) you will trigger installation. You can trigger as many installations as you want, they all be installed one-by-one. 
```lua
require("telescope").extensions.go_get.packages_search()
```

You can map the action for quicker usage like this

```lua
vim.keymap.set("n", "<Leader>gog", function()
  require("telescope").extensions.go_get.packages_search()
end, { desc = "[Go] [G]et packages" })
```


## Roadmap

There are things I would like to implement (but not limited to this):
- [x] Cleanup URL packages: ignore all forks
- [x] Validate new & old URL package: do not store URLs that cannot be installed anymore
- [x] Integrate package preview (in progress)
- [ ] Test integration with others plugin managers


## Sponsorship
If you find this project helpful and derive value from it, please consider supporting its development. I appreciate your interest and ask for your support to improve and maintain this project to benefit many others like you.

## Support and Q&A 

If you have any suggestion or improvement - feel free to open an issue or submit a PR.
If you would like to discuss - I invite you to my [Discord server](https://discord.gg/QeVvfvFfb6)

