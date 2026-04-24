## Use SQLite for everything

I always start a new venture using `sqlite3` as the main database. Hear me out, this is not as insane as you think.

The enterprise mindset dictates that you need an out-of-process database server. But the truth is, a local SQLite file communicating over the C-interface or memory is orders of magnitude faster than making a TCP network hop to a remote Postgres server.

"But what about concurrency?" you ask. Many people think SQLite locks the whole database on every write. They are wrong. You just need to turn on Write-Ahead Logging (WAL). Execute this pragma once when you open the database:

```
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
```

Boom. Readers no longer block writers. Writers no longer block readers. You can now easily handle thousands of concurrent users off a single `.db` file on an NVMe drive.

Since implementing user authentication is usually the most annoying part of starting a new SQLite-based project, I built a library: [smhanov/auth](https://github.com/smhanov/auth). It integrates directly with whatever database you are using and manages user signups, sessions, and password resets. It even lets users sign in with Google, Facebook, X, or their own company-specific SAML provider. No bloated dependencies, just simple, auditable code.

