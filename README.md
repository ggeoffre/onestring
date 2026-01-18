// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

One String, Many Prompts: A Hands-On Workshop in AI Code Generation
+++

What if one tiny JSON record could unlock an entire pipeline of working code? In this fast-paced, hands-on 2-hour workshop, we’ll use a single JSON string to generate cross-platform solutions with AI — spanning microcontrollers, APIs, and databases. Attendees won’t just see AI-generated code; they’ll actively craft prompts, refine them, and evaluate the results in real time. Along the way, we’ll reveal a powerful feedback loop: every prompt is captured, then fed back into AI for critique, showing how better organization and phrasing yield stronger outcomes. Whether you care about IoT tinkering, backend APIs, or data pipelines, you’ll leave with practical techniques for writing and refining prompts that turn AI into a productive coding partner. Whether you’re into IoT tinkering, backend APIs, or data engineering, this fast-paced session will show you how AI can accelerate development across the stack — starting from one simple string. 

*NOTE: Access the presentation slides here:*
<br><a href="https://www.slideshare.net/slideshow/one-string-many-prompts-ai-code-generation/285297740">
  <img src="https://image.slidesharecdn.com/onestringmanyprompts-260116182346-5971bd14/75/One-String-Many-Prompts-AI-Code-Generation-1-2048.jpg" width="400" alt="OneString Slide Presentation">
</a>

## TL;DR

**OneString is a hands-on method for making AI code generation repeatable.**  
By anchoring every exercise to a single, immutable JSON record, OneString teaches developers how to reduce variance in AI-generated code through clear intent, explicit constraints, and disciplined prompts. The focus isn’t on flashy demos — it’s on reliably reproducing working code across languages, frameworks, databases, and platforms using AI as a controlled engineering tool, not a conversational partner.

---

## Why OneString Works

**OneString works because it treats AI code generation as an engineering system, not a conversation.**

At its core is a single, stable JSON contract:

```json
{
  "recorded": 1768570200,
  "location": "den",
  "sensor": "bmp280",
  "measurement": "temperature",
  "units": "C",
  "value": 22.3
}
```

This constraint is intentional. It’s what makes everything else predictable.

---

### Constraints Reduce Variance

Large Language Models are probabilistic. Without constraints, they produce many “reasonable” answers — most of which drift from your intent.

OneString reduces variance by:
- Freezing the data schema
- Keeping field names and types constant
- Reusing the same payload across every stack

This mirrors real software engineering: stable interfaces matter more than clever implementations.

---

### Prompting Is System Configuration

OneString reframes AI concepts using engineering language:

| AI Term | Engineering Meaning |
|------|---------------------|
| Prompt | Intent expression |
| Tokens | Scarcity constraint |
| Output | Feedback surface |
| Retrieved context | Externalized knowledge |
| LLM | Probabilistic execution engine |

Prompting isn’t conversation — it’s configuration.  
Clear intent and explicit constraints produce reliable outcomes.

---

### Repeatability Over Novelty

OneString isn’t about generating impressive one-off results.

Every exercise follows the same loop:
1. Generate working code  
2. Run it  
3. Extract the code  
4. Turn it back into a prompt  
5. Refine the prompt to narrow variance  

The goal is simple: **reproduce the same working code on demand**.

---

### Instruction → Intent → Composition

The tasks are deliberately sequenced:
- **Data** → procedural correctness  
- **API** → interface contracts  
- **Access** → design intent (Inversion of Control)  
- **Legos** → system composition  

By the final exercise, the system assembles cleanly because each piece is already reliable in isolation.

---

### Scope Is Constrained on Purpose

The primary path starts with Python, Flask, and MongoDB — not because they’re “best,” but because they:
- Minimize boilerplate
- Reduce hallucinations
- Keep focus on prompt quality instead of tooling friction

Complexity comes later, once the mental model is solid.

---

### The Model Transfers

OneString doesn’t teach a language, framework, or AI tool.

It teaches how to:
- Express intent precisely
- Constrain probabilistic systems
- Design stable interfaces
- Build complex systems from small, reliable parts

Those skills transfer to any LLM, any stack, and any domain.

---

**OneString doesn’t try to make AI smarter.  
It teaches developers how to be more precise.**
