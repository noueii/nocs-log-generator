# CLAUDE.md - CS2 Log Generator Project Configuration

## Project Overview

This is a CS2 HTTP Log Generator project built by multiple specialized AI agents working in parallel. The project aims to generate realistic Counter-Strike 2 match logs for testing and analysis purposes.

**Current Phase**: MVP Development (Non-Production)
**Goal**: Create a working prototype that can generate CS2 logs and parse demo files

## Agent Orchestration System

### Available Agents

All agents are defined in `.claude/agents/` directory:
- **orchestrator** - Project coordination and task management
- **backend-engineer** - Go backend development
- **frontend-engineer** - React/TypeScript frontend
- **ui-designer** - shadcn/ui and Tailwind CSS implementation
- **test-engineer** - Testing and validation
- **devops-engineer** - Docker and deployment setup

### Agent Awareness Rules

1. **Context Awareness**: All agents MUST read and understand:
   - This CLAUDE.md file
   - PRD.md for project requirements
   - `.claude/tasks.md` for current task assignments
   - Other agents' work in their respective directories

2. **Parallel Work**: Agents can work simultaneously on:
   - Different modules (backend/frontend)
   - Independent features
   - Non-conflicting files

3. **Coordination Required**: Agents MUST coordinate when:
   - Modifying shared interfaces or contracts
   - Changing API specifications
   - Updating database schemas
   - Modifying configuration files

4. **Session Resumption**: Each agent MUST:
   - Check `.claude/tasks.md` for their assigned tasks
   - Review recent commits before starting work
   - Update task status when starting/completing work
   - Document any blockers or dependencies

## Coding Standards

### General Principles

1. **DRY (Don't Repeat Yourself)**
   - No code duplication across modules
   - Create shared utilities in appropriate directories
   - Use constants for repeated values

2. **SOLID Principles**
   - Single Responsibility: Each module/component does one thing
   - Open/Closed: Extensible without modification
   - Liskov Substitution: Interfaces must be properly implemented
   - Interface Segregation: Small, focused interfaces
   - Dependency Inversion: Depend on abstractions, not concretions

### Naming Conventions

#### Go Backend
```go
// Packages: lowercase, single word
package generator

// Files: snake_case
event_generator.go
match_config.go

// Types: PascalCase
type MatchConfig struct {}
type EventGenerator interface {}

// Functions/Methods: PascalCase for exported, camelCase for private
func GenerateMatch() {}
func parseConfig() {}

// Constants: PascalCase
const MaxRounds = 24
const DefaultTickRate = 64

// Variables: camelCase
var matchState = &MatchState{}
var currentRound int
```

#### React/TypeScript Frontend
```typescript
// Files: PascalCase for components, camelCase for utilities
MatchGenerator.tsx
PlayerCard.tsx
matchUtils.ts
apiClient.ts

// Components: PascalCase
export function MatchGenerator() {}
export const PlayerCard: FC<Props> = () => {}

// Hooks: use prefix
export function useMatch() {}
export function useWebSocket() {}

// Types/Interfaces: PascalCase with I/T prefix for clarity
interface IMatchConfig {}
type TPlayerRole = 'entry' | 'awp' | 'support'

// Constants: UPPER_SNAKE_CASE
const MAX_PLAYERS = 10
const API_BASE_URL = '/api/v1'

// Variables/Functions: camelCase
const matchData = {}
function calculateEconomy() {}
```

### Project Structure

```
/backend
  /cmd/server      - Main application entry
  /pkg
    /api          - HTTP handlers
    /generator    - Log generation logic
    /parser       - Demo parsing logic
    /models       - Data structures
    /utils        - Shared utilities

/frontend
  /src
    /components   - React components
    /pages        - Page components  
    /hooks        - Custom hooks
    /lib          - Utilities
    /services     - API clients
    /store        - State management
    /types        - TypeScript types

/.claude          - AI agent configuration
  /agents         - Agent definitions
  tasks.md        - Task tracking
```

## MCP Server Usage

### Available MCP Servers

Agents should leverage these MCP servers appropriately:

1. **filesystem** - File operations
   - Use for reading/writing project files
   - Navigate project structure

2. **github** - GitHub operations
   - Create issues for bugs
   - Reference PRs and commits
   - Check repository status

3. **sequential-thinking** - Complex problem solving
   - Use for architectural decisions
   - Algorithm design
   - Complex debugging

4. **memory** - Persistent memory
   - Store context between sessions
   - Remember design decisions
   - Track completed tasks

5. **shadcn-ui** - UI component generation
   - Generate shadcn/ui components
   - Follow design system

### MCP Usage Guidelines

```markdown
## When to use MCP servers:

### filesystem
- Reading multiple related files
- Batch file operations
- Project-wide searches

### github  
- Checking issue status
- Reviewing PR feedback
- Understanding commit history

### sequential-thinking
- Designing complex algorithms
- Solving architectural problems
- Planning implementation strategy

### memory
- Storing session context
- Remembering design decisions
- Tracking long-term progress

### shadcn-ui
- Generating new UI components
- Ensuring consistent design
- Following shadcn patterns
```

## Git Workflow

### Commit Standards

All commits MUST follow this format:
```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/config changes

Example:
```
feat(backend): add match generation endpoint

- Implement POST /api/generate endpoint
- Add match configuration validation
- Support MR12 and MR15 formats

Task: #TASK-001
```

### Branch Strategy

For MVP development, we work directly on main branch with:
- Frequent, small commits
- Clear commit messages
- Task references in commits

## Development Phases

### Phase 1: Foundation (Current)
- [ ] Project structure setup
- [ ] Basic backend API
- [ ] Frontend scaffold with shadcn/ui
- [ ] Docker configuration

### Phase 2: Core Features
- [ ] Match generation logic
- [ ] Basic UI for configuration
- [ ] Log output formatting
- [ ] Simple demo parsing

### Phase 3: Enhancement
- [ ] Advanced configuration options
- [ ] Real-time log streaming
- [ ] Demo file upload UI
- [ ] Basic statistics

### Human Validation Gates

Each phase requires human validation before proceeding:
1. Code review of implementation
2. Manual testing of features
3. Approval to proceed to next phase

## Task Management

Tasks are tracked in `.claude/tasks.md` with format:
```markdown
## TASK-XXX: Task Title
- **Assigned**: agent-name
- **Status**: pending|in-progress|completed|blocked
- **Priority**: high|medium|low
- **Dependencies**: TASK-YYY, TASK-ZZZ
- **Description**: Clear task description
```

## Quality Checklist

Before marking any task complete, ensure:

- [ ] Code follows naming conventions
- [ ] No duplicate code exists
- [ ] Proper error handling added
- [ ] Basic tests written (if applicable)
- [ ] Documentation updated
- [ ] Task status updated in tasks.md
- [ ] Commit follows standards

## Important Notes

1. **MVP Focus**: We're building a working prototype, not production code
2. **Iterative Development**: Start simple, enhance gradually
3. **Human Testing**: Each feature tested by humans before moving forward
4. **Clear Communication**: Document blockers and dependencies
5. **Context Preservation**: Always update task status for session resumption

## Agent Communication Protocol

When agents need to communicate:
1. Update task status in `.claude/tasks.md`
2. Leave clear comments in code for other agents
3. Document API contracts in code
4. Create interface definitions before implementation
5. Flag breaking changes in commit messages

## Error Handling Standards

### Backend (Go)
```go
// Always return errors, don't panic
func GenerateMatch(config *MatchConfig) (*Match, error) {
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    // ...
}
```

### Frontend (React)
```typescript
// Use try-catch with proper error boundaries
try {
    const match = await generateMatch(config)
} catch (error) {
    console.error('Match generation failed:', error)
    toast.error('Failed to generate match')
}
```

## Testing Approach

For MVP phase:
- Focus on happy path testing
- Basic validation tests
- Manual testing by humans
- No need for 100% coverage yet

## Remember

- Check `.claude/tasks.md` before starting work
- Read other agents' recent work
- Follow the established patterns
- Ask orchestrator if unclear about task
- Update progress regularly
- Commit frequently with clear messages