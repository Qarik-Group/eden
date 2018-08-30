# Improvements

- New `--json` flag for the `catalog` and `services` commands, to
  make it easier to use via jq (without spruce).

- The `bind` and `unbind` commands now support an `SB_BINDING`
  environment variable.  Additionally, you can now specify your
  binding ID ahead of time for the `bind` command.
