# Sunrise

A 'light' which when turned on runs a simple animation to simulate a sunrise with your Hue lights.

https://streamable.com/wb1rr

## Configuration

- Duration: to change the duration of the sunrise, change the value of the `EndTransitionTime` variable
  (multiples of 100ms).
- Lights affected: This is currently hardcoded to apply the sunrise to the first group defined, whatever that is.
  To change this behavior, you will have to edit the code [here](https://github.com/ragurney/hue-config/blob/b982267ba1168cb99417ecf2a58fb65dce85c6b8/animations/sunrise/sunrise.go#L69).
