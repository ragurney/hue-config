# Sunrise

A 'light' which when turned on runs a simple animation to simulate a sunrise with your Hue lights.

https://streamable.com/wb1rr

## Configuration

- Duration: to change the duration of the sunrise, change the value of the `EndTransitionTime` variable
  (multiples of 100ms).
- Lights affected: This is currently hardcoded to apply the sunrise to the first group defined, whatever that is.
  To change this behavior, you will have to edit the code under the comment:
  ```
  // Sends sunrise command to the first group, whatever that is. Can customize this to your own needs.
  ```
