# RSS Builder

Build your own custom RSS feeds based on any website you like.

## Overview

RSS Builder is a tool that lets you create custom RSS feeds from virtually any website by specifying your own
CSS/GoQuery selectors.  
The idea is straightforward: you tell the system how to find the content you care about—titles, article URLs, images, or full
content—and it will automatically collect it and store it in a database.

You can then generate an RSS feed from the latest collected content, giving you a way to track websites that do not
provide an RSS feed natively.

> [!WARNING] This project is still in a very early development phase. Features may be incomplete, unstable, or
> subject to change.
> This is a personal project to learn Go, so don't expect anything stable.

---

## Features

- Define custom scraping targets using your own CSS selectors
- Track titles, article URLs, images, and full content
- Store scraping results in a PostgreSQL database
- Generate an RSS feed from the latest scraping process
- Designed to be extensible for multiple targets and feeds

## Getting Started

### Prerequisites

- Go 1.25+ (or newer)
- PostgreSQL 12+ (for storing targets, processes, and results)

Notes

* The project is still in early development: features and APIs may change frequently.
* Currently, no UI, automation, or scheduling – everything is triggered manually.
* Designed to be extensible for more features like multiple targets, scheduling, or advanced feed customization.
