# go-zed-test-toggle Examples

This directory contains examples of how to use `go-zed-test-toggle` with different Ruby project structures.

## Example Project Structures

### 1. Ruby Gem with RSpec

```
my-gem/
├── .rspec
├── my-gem.gemspec
├── lib/
│   ├── my_gem.rb
│   └── my_gem/
│       ├── user.rb
│       └── post.rb
└── spec/
    ├── spec_helper.rb
    ├── my_gem_spec.rb
    └── lib/
        └── my_gem/
            ├── user_spec.rb
            └── post_spec.rb
```

**Toggle examples:**
- `lib/my_gem/user.rb` ↔ `spec/lib/my_gem/user_spec.rb`
- `lib/my_gem.rb` ↔ `spec/lib/my_gem_spec.rb`

### 2. Rails Application with RSpec

```
rails-app/
├── .rspec
├── app/
│   ├── models/
│   │   └── user.rb
│   ├── controllers/
│   │   └── users_controller.rb
│   └── services/
│       └── user_service.rb
├── lib/
│   └── tasks/
│       └── users.rake
└── spec/
    ├── spec_helper.rb
    ├── models/
    │   └── user_spec.rb
    ├── controllers/
    │   └── users_controller_spec.rb
    └── services/
        └── user_service_spec.rb
```

**Toggle examples:**
- `app/models/user.rb` ↔ `spec/models/user_spec.rb`
- `app/controllers/users_controller.rb` ↔ `spec/controllers/users_controller_spec.rb`
- `app/services/user_service.rb` ↔ `spec/services/user_service_spec.rb`

### 3. Ruby Project with Minitest

```
minitest-project/
├── lib/
│   └── calculator.rb
└── test/
    ├── test_helper.rb
    ├── calculator_test.rb
    └── test_calculator.rb  # Alternative naming
```

**Toggle examples:**
- `lib/calculator.rb` ↔ `test/calculator_test.rb`
- `lib/calculator.rb` ↔ `test/test_calculator.rb`

## Zed Configuration

### Basic Task Configuration

Add to your Zed tasks (`~/.config/zed/tasks.json` or project `.zed/tasks.json`):

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

Add to your Zed keybindings (`~/.config/zed/keymap.json`):

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

## Usage Examples

### Command Line Usage

```bash
# From a Ruby gem project
go-zed-test-toggle lookup -p "lib/my_gem/user.rb" -r "/path/to/my-gem"
# Opens: /path/to/my-gem/spec/lib/my_gem/user_spec.rb

# From a Rails project
go-zed-test-toggle lookup -p "app/models/user.rb" -r "/path/to/rails-app"
# Opens: /path/to/rails-app/spec/models/user_spec.rb

# Toggle back from test to source
go-zed-test-toggle lookup -p "spec/models/user_spec.rb" -r "/path/to/rails-app"
# Opens: /path/to/rails-app/app/models/user.rb
```

### Within Zed

1. Open any Ruby source or test file
2. Press `Ctrl+Shift+T` (or your configured keybinding)
3. The corresponding test or source file will open

## Special Cases

### Gem Root Files

For gem projects, files in the root directory are mapped to `spec/` directory:

```
my_gem.rb → spec/my_gem_spec.rb
```

### Nested lib Directories

For projects with nested lib structures in tests:

```
lib/my_gem/core/engine.rb → spec/lib/my_gem/core/engine_spec.rb
```

### Multiple Test Naming Conventions

The tool recognizes both Minitest naming conventions:
- `*_test.rb` (e.g., `user_test.rb`)
- `test_*.rb` (e.g., `test_user.rb`)

## Troubleshooting

### File Not Found

If the alternate file is not found, the tool exits silently. Common reasons:

1. **Test file doesn't exist yet**: Create the test file manually
2. **Non-standard directory structure**: The tool expects conventional Ruby project structures
3. **Wrong project root**: Ensure `-r` points to the project root, not a subdirectory

### Detection Issues

The tool detects project type by looking for:

**RSpec indicators:**
- `.rspec` file in project root
- `spec/spec_helper.rb`
- Any `spec_helper.rb` in the project

**Gem indicators:**
- `*.gemspec` file in project root

### Performance

The Go implementation is significantly faster than the Ruby version:
- Instant startup (no Ruby interpreter overhead)
- Efficient file system operations
- Single binary with no dependencies

## Integration Tips

### VS Code Integration

While designed for Zed, you can use this with VS Code by creating a task:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Toggle Test",
      "type": "shell",
      "command": "go-zed-test-toggle",
      "args": [
        "lookup",
        "-p",
        "${relativeFile}",
        "-r",
        "${workspaceFolder}"
      ]
    }
  ]
}
```

### Shell Alias

For quick command-line usage:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias toggle-test='go-zed-test-toggle lookup -p "$(pwd)/${1#./}" -r "$(git rev-parse --show-toplevel 2>/dev/null || pwd)"'

# Usage
toggle-test lib/user.rb
```

### Git Hooks

Use in pre-commit hooks to ensure tests exist:

```bash
#!/bin/bash
# .git/hooks/pre-commit

for file in $(git diff --cached --name-only | grep -E '\.rb$'); do
  if [[ ! "$file" =~ _(spec|test)\.rb$ ]]; then
    alternate=$(go-zed-test-toggle lookup -p "$file" -r "$(pwd)" 2>/dev/null)
    if [ -z "$alternate" ]; then
      echo "Warning: No test found for $file"
    fi
  fi
done
```
