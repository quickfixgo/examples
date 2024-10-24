
# Load Test for TradeClient

## Overview
This load test evaluates the performance of the TradeClient by simulating the submission of multiple orders in a specified time frame. The primary focus is to measure the success/failure percentage of the orders processed, while performance analysis will be derived from log files.

## Types of Tests Running
- **Load Test**: Simulates the submission of multiple orders to evaluate the success and failure rates.
- **Performance Metrics**: Uses existing logs to analyze response times and latencies for orders processed.

## How to Run the Tests
1. **Build the TradeClient**:
   - Use one of the following commands:
     ```make build```
     or
     ```make clean build```

2. **Run the TradeClient**:
   - Execute the following command:
     ```./bin/qf tradeclient```

3. **Select Load Test**:
   - You will be prompted with the following options:
        1) Enter Order
        2) Cancel Order
        3) Request Market Data
        4) Run Load Test
        5) Quit
   - Choose **4** to initiate the load test.

## Outputs
- The results of the load test, including success and failure rates, will be logged to `output.log`.
