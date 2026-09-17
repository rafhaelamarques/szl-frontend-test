# AI-Assisted Development

This document records the prompts used to guide AI-assisted planning and implementation for this project. The prompts have been lightly edited for clarity and consistency without changing their intended instructions.

## Planning Prompt

Read `assignment.txt` and produce a concrete implementation plan for the greenfield full-stack calculator application in this repository.

### Developer context

I am a mobile engineer with professional experience in Flutter, Kotlin, Swift, REST APIs, Java and SpringBoot. I have limited production experience with TypeScript, Go and Docker.

When relevant, briefly relate unfamiliar concepts to these technologies. Keep the recommendations idiomatic for React and Go.

### Scope and quality expectations

Treat `assignment.txt` as the source of truth. Clearly distinguish mandatory requirements from optional features and recommend a scope that can realistically be completed within the expected two-to-four-hour timeframe.

Favor:

- Idiomatic React and Go
- Clear, readable, maintainable code
- A simple architecture appropriate for a small calculator
- Practical testing and error handling
- Decisions that are easy to explain during a technical review

Avoid unnecessary frameworks, services, abstractions or infrastructure.

### Expected output

Provide an implementation plan only. Do not create or modify any files.

Structure the plan as follows:

1. Requirements and likely evaluation criteria
2. Complexity and main risks
3. Recommended repository structure and technical approach
4. API contract, validation, and error handling
5. Testing strategy
6. Docker and local development approach
7. Ordered implementation steps
8. Explicit assumptions and reviewer-visible decisions
9. Definition of done

Keep the plan concrete and actionable. When multiple approaches are viable, recommend one and briefly explain why it is appropriate for this assignment.

## Implementation Prompt

Implement the approved plan in this repository.

Follow the existing plan and project requirements. Keep the implementation proportional to the assignment's scope and time constraints. Prioritize correctness, clarity, maintainability, test coverage, and a straightforward developer experience.

Do not introduce additional frameworks, services, or abstractions unless they are necessary to satisfy the requirements.

Before completing the work:

- Run the relevant frontend and backend tests
- Run formatting and static checks where applicable
- Verify the application locally or through the documented Docker workflow
- Update the README with setup instructions, API examples, design decisions, and assumptions
- Report any remaining limitations or blockers clearly

Do not modify the planning document.

## Decisions Made During Planning

- Include all optional calculator operations: exponentiation, square root, and percentage.
- Use Vite, React, and TypeScript for the frontend.
- Use Go's standard library for the REST backend.
- Expose the calculation endpoint as `POST /api/v1/calculate`.
- Use table-driven tests for backend calculation logic.
- Use Vitest and Testing Library for frontend tests.
- Use Docker Compose with an Nginx `/api` proxy.
- Keep the README in English.
- Defer Git remote configuration until the implementation is complete.