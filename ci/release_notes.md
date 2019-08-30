*  Consider service plan `free` as optional (#20)

    As per specification `plan.free` is not a required parameter and entirely skipped in json serialization in the brokerapi implementation. (see omitempty).
    Eden will panic while trying to dereference a nil pointer.

* `eden` release built with Go 1.12