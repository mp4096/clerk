# Clerk

[![Build Status](https://travis-ci.org/mp4096/clerk.svg?branch=master)](https://travis-ci.org/mp4096/clerk)
[![Go Report Card](https://goreportcard.com/badge/github.com/mp4096/clerk)](https://goreportcard.com/report/github.com/mp4096/clerk)

Clerk sends your Markdown notes via email.
I use it for distributing meeting transcripts.

## Security

`clerk` uses the `net/smtp` package,
[which uses TLS if possible](https://golang.org/pkg/net/smtp/#SendMail).
Still, just to be safe, I explicitly discourage using `clerk` for mission-critical information.

## Installation

### From source on Linux

```
$ make install
```

### Binaries

You can get them from GitHub releases.

## Usage example

Suppose you've written a note in Markdown and saved it to `2017-12-04_note.md`.
Your configuration (see below) is defined in `jane.clerk.yml`.
You can now preview this note rendered as an HTML email:

```
$ clerk distribute -m 2017-12-04_note.md -c jane.clerk.yml
Hello, Jane Doe
Send flag not set: opening preview in "chromium-browser"
```

And send the email by adding the send flag `-s`:

```
$ clerk distribute -m 2017-12-04_note.md -c jane.clerk.yml -s
Hello, Jane Doe
Will send to [john.doe@abc.com max.mustermann@def.de jane.doe@xyz.com]
Please enter your credentials for "smtpserver.xyz.com"
Login: janedoe
Password:
```

Since this email will not appear in your provider's `Sent` folder,
`clerk` will send you a BCC copy.

Since `clerk` was designed for meeting transcripts,
you can also use the `clerk approve` command to send an email with approval request to your boss.

For help, call `clerk help` or `clerk <command> -h`.

Important: The filename _must_ begin with a valid ISO date.
It will be parsed and filled in as a context variable `{{date}}`.

## How to configure

Here's an example config file:

```yaml
email_server:
    hostname: smtpserver.xyz.com
    port: 587

author:
    name: Jane Doe
    email: jane.doe@xyz.com
    notice: |
        <i>This email was sent with <tt>clerk</tt>
        (<a href="https://github.com/mp4096/clerk" target="_top">get it on GitHub</a>).
        </i>
    browser: chromium-browser

approve_list:
    emails:
        - john.doe@abc.com
    subject: Approval request
    salutation: Hi John,
    text: |
        here are my notes from {{date}}. What do you think?

distribute_list:
    emails:
        - john.doe@abc.com
        - max.mustermann@def.de
    subject: My thoughts on {{date}}
    salutation: Hello everyone,
    text: |
        here are my notes from {{date}}.
```
