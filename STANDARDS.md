These are requirements for code using the orc framework.

## Security properties

- The values of flags must never be secret.
- The values of environment variables beginning with "ORC\_" (ORC
  followed by an underscore) must never be secret.
- Secret values should be passed as files, with a flag to specify
  the file name.

## Expectations

- orc is a hobby project and primarily intended for the author's
  own usage. Currently, no support can be expected.
- Because of the above lack of support, don't rely on orc for
  anything security-critical or anywhere where the consequences
  of bugs could be severe.

## Dependencies and integrations

Orc code uses, and assumes your application uses:

- `cobra` to manage its flags and subcommands.
- `logrus` for logging.
- Prometheus for monitoring.

## Deviations from normal coding style

- Orc code uses global state in modules designed to be used as singletons.
- Orc code declares flags in module libraries.
- Orc code may log from module library code.
- Orc code may exit the program (logrus.Fatal) from module library code
  upon unrecoverable errors on initialization failures.
- In general, orc favours usability, debuggability, and ease of
  programming over heavy performance optimizations.

## Writing orc modules

- In implementations of Orc modules, references to values from other
  modules shall be done through methods in a module object, and
  shall be confined to a file called orc.go or module.go.
- By convention, when a module is expected to be used as a singleton
  module which does not require any configuration, it should declare
  its module type as "Module" and create a ready-to-use instance of it
  as "M".
