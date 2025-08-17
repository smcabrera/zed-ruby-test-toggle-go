# Go Zed Test Toggle

A Go implementation of the `zed-test-toggle` gem, designed to allow developers to quickly switch between source files and their corresponding test files in the Zed editor. This tool supports both RSpec (`spec` directory) and Minitest/Test::Unit (`test` directory) conventions.

## Features

- Fast switching between source files and test files
- Supports Ruby/Rails project structures
- Auto-detects RSpec vs Minitest/Test::Unit conventions
- Works with gem projects and Rails applications
- Lightweight and fast Go implementation

## Installation

### From Source

```bash
git clone https://github.com/stephen/go-zed-test-toggle.git
cd go-zed-test-toggle
go build -o go-zed-test-toggle
```

Move the binary to a location in your PATH:

```bash
sudo mv go-zed-test-toggle /usr/local/bin/
# or
mv go-zed-test-toggle ~/bin/  # if ~/bin is in your PATH
```

### Using Go Install

```bash
go install github.com/stephen/go-zed-test-toggle@latest
```

## Usage

This tool is designed to be called from a Zed task. Add the following to your Zed tasks configuration:

```json
[
  {
    "label": "Toggle Test and Target",
    "command": "go-zed-test-toggle",
    "args": [
      "lookup",
      "-p",
      "\"$ZED_RELATIVE_FILE\"",
      "-r",
      "\"$ZED_WORKTREE_ROOT\""
    ],
    "hide": "always",
    "allow_concurrent_runs": false,
    "use_new_terminal": false,
    "reveal": "never"
  }
]
```

### Keybinding Configuration

To use the tool effectively, configure a keybinding. Add this to your Zed keybindings:

```json
[
  {
    "bindings": {
      "ctrl-shift-t": [
        "task::Spawn",
        {
          "task_name": "Toggle Test and Target",
          "reevaluate_context": true
        }
      ]
    }
  }
]
```

**Note:** The `reevaluate_context: true` option is crucial. Without it, the environment variables won't be refreshed, and you'll keep jumping to the same file.

## How It Works

The tool uses the following logic to find alternate files:

### Project Detection

1. **Gem Detection**: Looks for `*.gemspec` files to determine if it's a gem project
2. **Test Framework Detection**: 
   - RSpec: Looks for `spec/spec_helper.rb`, `.rspec`, or any `spec_helper.rb` file
   - Minitest/Test::Unit: Default if RSpec indicators aren't found

### Path Mapping

#### For Gem Projects
- Source paths: `lib/`, root directory
- Test paths: `spec/` or `test/`, `spec/lib/` or `test/lib/`

#### For Rails/Regular Projects
- Source paths: `app/`, `lib/`
- Test paths: `spec/` or `test/`, `spec/lib/` or `test/lib/`

### File Name Conventions

#### RSpec
- Test files end with `_spec.rb`
- Source file `lib/user.rb` → Test file `spec/lib/user_spec.rb`

#### Minitest/Test::Unit
- Test files end with `_test.rb` or match `test_*.rb`
- Source file `lib/user.rb` → Test file `test/lib/user_test.rb`

## Examples

Given a project structure:
```
myproject/
├── lib/
│   └── models/
│       └── user.rb
├── spec/
│   └── lib/
│       └── models/
│           └── user_spec.rb
└── .rspec
```

- When editing `lib/models/user.rb`, pressing `ctrl-shift-t` opens `spec/lib/models/user_spec.rb`
- When editing `spec/lib/models/user_spec.rb`, pressing `ctrl-shift-t` opens `lib/models/user.rb`

## Differences from Ruby Implementation

This Go implementation aims to be functionally equivalent to the original Ruby gem with the following characteristics:

- **Performance**: Faster startup time due to compiled binary vs Ruby interpreter
- **No Dependencies**: Single binary with no runtime dependencies
- **Same Logic**: Maintains the same file detection and mapping logic as the Ruby version

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the same license as the original zed-test-toggle gem.

## Acknowledgments

This is a Go port of the [zed-test-toggle](https://github.com/MoskitoHero/zed-test-toggle) Ruby gem by MoskitoHero (Cédric Delalande).