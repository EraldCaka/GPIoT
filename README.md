# GPIoT

GPIoT is a service that utilizes the MQTT protocol to enable real-time data transfer between a serverand IoT devices like Raspberry Pi GPIO pins.It supports both digital and analog pins for various IoTprojects.

---

⚠️ **Caution**: GPIO pins are simulated in this implementation. For real GPIO connectivity, you may need to integrate with specific hardware libraries suitable for your device.

---


## Features

- **Real-time Communication**: Leverages MQTT for efficient data exchange.
- **Support for Digital and Analog Pins**: Read and write states for GPIO pins.
- **Easy Integration**: Designed to work seamlessly with Raspberry Pi and other devices.
- **Configurable**: Allows customization of pins, modes, and topics via configuration files.

---

## How It Works

1. **Initialization**:
   - The application initializes GPIO pins based on the configuration file.
   - Pins can be configured as digital (input/output) or analog (input/output/both).

2. **MQTT Communication**:
   - The service subscribes to specified topics to listen for control messages.
   - Publishes pin states to response topics for real-time monitoring.

3. **Event Handling**:
   - Monitors GPIO pins and triggers events based on their states.
   - Reads and writes analog values from/to files for persistent storage.

---

## Installation

### Prerequisites

- **Hardware**:
  - Raspberry Pi (or compatible IoT device)
  - GPIO Pins with required peripherals

- **Software**:
  - Go (Golang) installed
  - MQTT broker (e.g., Mosquitto)

### Steps

1. **Clone the Repository**:
   ```bash
   $ git clone https://github.com/EraldCaka/GPIoT.git
   $ cd GPIoT
   ```

2. **Install Dependencies**:
   ```bash
   $ go mod tidy
   ```

3. **Configure the Application**:
   - Edit the `config.json` file to define your MQTT settings and GPIO pins.

4. **Run the broker**:
    ```bash
    $ make up
    ```

5. **Run the Application**:
   ```bash
   $ go run main.go
   ```
6. **Run Integration Tests**:

  After running the broker and the application you can run these integration
  tests to change analog or digital pins values.

  - Analog Pin Integration:
  ```bash
  $ make analog
  ```
  - Digital Pin Integration:
  ```bash
  $ make digital
  ```

---

## Usage

- **Listening for Messages**:
  - The application subscribes to MQTT topics to receive control messages.

- **Sending Control Messages**:
  - Publish messages to configured topics to control GPIO pins.

- **Monitoring**:
  - Pin states are published to response topics, allowing real-time tracking.

---

## Client

A third party client can be connected with the broker to monitor
the incoming messages and published ones. A good one would be ``MQTTX``.

### MQTTX
![mqttx](/public/mqttx.png)

## Example Configuration

```yaml
broker: tcp://localhost:1883
client_id: gpio-monitor
monitor_time: 5s
mqtt:
  qos: 0
  retain: false
gpio:
  digital-pins:
    - gpio-name: pir
      pin: 2
      mode: 1
      state: 0
      path: gpio-mock-store/pin2.txt
      topic: gpio/control/digital/2
  analog-pins:
    - gpio-name: fig
      pin: 4
      mode: 2
      state: 0.5
      path: gpio-mock-store/pin4.txt
      topic: gpio/control/analog/4
```
