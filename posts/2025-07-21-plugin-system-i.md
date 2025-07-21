---
title: "Plugin Systems Part I"
date: 2025-07-21
teaser: How hard can it be?
background: "assets/2025-07-21-background.png"
---

## Why do a plugin system?
Because it's cool! But there's so much more. 
There are the technical aspects. 
Plugins are the right hammer to certain problems (we'll get to that later). 
And there are business strategic, sometimes even political, aspects. 
Let's have quick dive into that first.

![image](../assets/2025-07-21-illustration-1.png)

## Case Study 1: VSCode
VSCode without it's plugin system (let's call them extensions and marketplace from now on to use their terminology), 
without extensions it's is just a dumb text editor wrapped around a browser. 
Extensions are where the magic comes into play. 
It's a clever lego system that provides you with powerful building blocks so you can easily create your own software factory. 
Or not and just use it as an editor to write your thesis in LaTex. 
Some big enterprise companies even use it to write VHDL and build hardware with it. 
And if you cobble too much of it together, it's gonna fall apart and become utterly slow and convoluted. 
Judging by the VSCodes IDE market share of ~74% ([Technology | 2024 Stack Overflow Developer Survey](https://survey.stackoverflow.co/2024/technology)), 
for most devs it's way more than a simple editro but comes close enough to a full fledged IDE so to them there's effectively no difference.

But there's also a not so obvious flip side to that success story. 
While the core of VSCode and many of the extensions are open source, 
Microsoft has put severe restrictions around. Marketplace extensions are _only_ licensed for "In-Scope Products and Services"
(Visual Studio, VS Code, GitHub Codespaces, Azure DevOps, etc.).
Third-party apps must not "install, import, reverse-engineer, scrape or spider" the store or its contents ([source](https://code.visualstudio.com/license?ref=ghuntley.com)).
So while it's true VSCode is open source, this effectively cuts everyone other than Microsoft out from re-using most of the flourished ecosystem around it.
That has consequences and some people are very unhappy about it to put it mildly. It's been a big culprit eg for GitPod, 
almost taking down their business when Microsoft/Github launched their own version of it, 
Github Workspaces (this rant is very much worth reading about this in detail: [Visual Studio Code is designed to fracture](https://ghuntley.com/fracture/)). 
They weren't allowed to offer people all the richness of the VSCode marketplace which made look things pretty empty. 
That again was the case for Google when they launched IDX, a VSCode derivative for the AI browser IDE. 
It again lacks the richness of the VSCode marketplace [idx announcement](https://developers.googleblog.com/en/introducing-project-idx-an-experiment-to-improve-full-stack-multiplatform-app-development/). 
And beside being cut out of the marketplace ecosystem aspect, it's even worse:
> Microsoft forks open-source communities by releasing Visual Studio Code extension updates that make their proprietary 
> offering the default once they have managed to capture enough adoption. 
> (Geoffrey Huntley, [Visual Studio Code is designed to fracture](https://ghuntley.com/fracture/))

See [Open Source VS Code Python Language Server Dies, Replaced by Proprietary Pylance -- Visual Studio Magazine](https://visualstudiomagazine.com/articles/2021/11/05/vscode-python-nov21.aspx?ref=ghuntley.com) for the implication on the Python LSP.

But that of course doesn't mean plugins or extension are pure evil. 
It also doesn't mean VSCode is evil. 
In contrary, VSCode is probably one of the best things that has happened to dev tooling in recent years. 
It has pushed boundaries on how we think about IDEs and has gifted the dev community a few useful standards, most popular LSP. 
It's the protocol that powers syntax highlighting and now is also used e.g. in Neovim ([Lsp - Neovim docs](https://neovim.io/doc/user/lsp.html)) or Zed ([Configuring Languages - Zed](https://zed.dev/docs/configuring-languages)). 
And even Jetbrains supports it now, at least to some extend ([Language Server Protocol (LSP) for Plugin Developers | The JetBrains Platform Blog](https://blog.jetbrains.com/platform/2023/07/lsp-for-plugin-developers/)). 
This is great to all developers. When some new language/framework/.. comes up we can expect it be supported across various IDEs much sooner. 
Also there's the [Development Container Specification](https://containers.dev/implementors/spec/) and lesser known the debug adapter protocol (DAP).

![image](../assets/2025-07-21-illustration-2.png)

## Case Study 2: Jetbrains
Jetbrains is again plugins everywhere but paired with a very different architecture and business strategy. 
The core platform [IntelliJ](https://github.com/JetBrains/intellij-community) (which also is their free/community Java IDE) same as VSCode is open source, 
however, it already comes as "batteries included" package with a fair amount of refactoring tools, syntax highlighting etc. 
By default, this sets it above the level of pure editors and positions it along professional IDEs. 
The technical approach to plugins could be described as more "integrated". 
It's based on PicoContainer, a lightweight Java dependency injection framework. 
De facto that means their IDE is a highly extendable platform with integration points everywhere and 
the official extension point browser currently lists 1700 of them ([Extensions | IntelliJ Platform Plugin SDK](https://plugins.jetbrains.com/docs/intellij/plugin-extensions.html))! 
And while Jetbrains is less known for being the inventor of a specific protocol or standard, 
that approach led them invented an entire new language Kotlin. 
They also are heavily engaged with the dev community but in a different way:

> We at JetBrains love going to conferences around the world. You can meet us at over 170 conferences each year! ([Conferences - JetBrains](https://www.jetbrains.com/company/conferences/))

Ie., they effectively are sending their devs around the world to pic up the latest dev trends and 
then because of their highly plug-able platform can easily adapt. 
Comparing with VSCode which is more like a lego system, the Jetbrains strategy is more like a slowly morphing platform
that constantly keeps adjusting and evolving to match the industry demands. 
That doesn't go without some limitations. 
If you keep carrying around what looks a lot like a heavy monolith, certain more radical quick jumps are harder. 
While we have seen a series of VSCode forks/derivatives spearheading the front of new AI driven dev tools
(Cursor, IDX, Windsurf, ...), there's way less noise in around people forking IntelliJ. 
And there also seems to some controversy about their own attempts ([JetBrains defends removal of negative reviews for unpopular AI Assistant | Hacker News](https://news.ycombinator.com/item?id=43850377)).

