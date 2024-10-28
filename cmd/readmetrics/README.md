# ReadMetrics

The `readmetrics` command is a tool for reading and processing metrics from a FIX log file. It calculates the latency and intervals between messages and logs them to a separate metrics file for analysis.

## Features

- Reads FIX log files and extracts relevant timestamps.
- Calculates latency and interval in milliseconds.
- Outputs metrics to a `metrics.log` file.

## Usage

To use the `readmetrics` command, run the following command in your terminal:

```
./bin/qf readmetrics <your_fix_log_file.log>
```

Replace `<your_fix_log_file.log>` with the path to your FIX log file.

## Output

The results will be written to a file named `metrics.log` in the current directory. The log will contain entries formatted as follows:

```
Message 1: Interval = X.XXXms, Latency = Y.YYYms
```

Where `X.XXX` is the interval between messages, and `Y.YYY` is the calculated latency.

## Example

Assuming you have a FIX log file named `example.log`, you can run:

```
./bin/qf readmetrics example.log
```

After execution, you can check the contents of `metrics.log`:

You should see output similar to:

```
Message 1: Interval = 1.000ms, Latency = 0.647ms
Message 2: Interval = 1.000ms, Latency = 0.685ms
...
```