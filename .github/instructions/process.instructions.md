# Process and Tooling Instructions

## Core Principle: Use The Proper Tools

**DO NOT bypass established tooling to achieve outcomes faster.**

Every tool in this project exists for a reason:

- **CLI tools**: Provide validation, consistency checks, proper metadata management
- **Code generators**: Ensure type safety, maintain generated code contracts
- **Build systems**: Manage dependencies, compilation, deployment

**WRONG**: Manually editing files to achieve the desired state faster
**RIGHT**: Using the proper tool even if it takes more steps

## Golden Rule: If a Tool Exists, Use It

Before writing custom code or manually editing files, ask:

1. Does a CLI tool exist for this? → Use it
2. Does a generator exist for this? → Use it
3. Is this a generated file? → NEVER edit directly, use the generator
4. Have we solved this before? → Use the existing pattern

## Examples of Proper Tool Usage

### GraphQL Schema Changes

**WRONG**: Edit `generated.go` or `*.resolvers.generated.go` directly
**RIGHT**:

1. Edit `.graphql` schema files
2. Run `go run github.com/99designs/gqlgen generate`
3. Implement panic stubs in generated resolver files

### Kanbn Board Updates

**WRONG**: Edit `.kanbn/index.md` directly
**RIGHT**: Use `kanbn move <task-id> -c "Column"`

### Database Schema Changes (future)

**WRONG**: Write SQL directly or edit migration files
**RIGHT**: Use migration tool CLI to generate migrations

## When in Doubt

1. Check if a CLI tool exists: `which <tool>`, `<tool> help`
2. Check project documentation in `docs/` and `.github/instructions/`
3. Ask before reinventing

## Credit Efficiency

Every tool call costs credits. Before executing:

1. **Read available output carefully** - the answer is often already there
2. **Use correct syntax** - check help output format before running commands
3. **Don't trial-and-error** - understand the tool first
4. **Batch operations** when possible

## DRY: Don't Repeat Yourself

This applies to work processes too:

- If we've established a correct way to do something, **use that way**
- Don't rediscover solutions we've already implemented
- Reference existing patterns and documentation
- When unsure, check docs/instructions first

## Checks and Balances

Never bypass validation:

- Use `gqlgen generate` instead of editing generated files
- Use package managers instead of manual downloads
- Use CLI tools instead of editing config files manually
- Let tools update metadata, timestamps, and cross-references

## Committing Code

When I say commit:

- Assume I mean everything. `git add -A && git commit -m "[meaning ful commemt about the commit]"`.
- If I say "commit X", I mean "commit X and all related changes needed to make X work correctly".

## Summary

**Use the right tool. Follow the process. Don't waste credits on shortcuts that bypass validation.**
