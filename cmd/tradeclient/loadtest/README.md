
# Load Test for TradeClient

## Overview
This load test evaluates the performance of the TradeClient by simulating the submission of multiple orders in a specified time frame. The primary focus is to measure the success/failure percentage of the orders processed, while performance analysis will be derived from log files.

## Types of Tests Running
- **Load Test**:  Simulates the submission of a high volume of orders at configurable rates.

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
      5) Read metrics
      6) Quit
   - Choose **4** to initiate the load test.

## Outputs
   - Log files will be generated in the tmp folder. To view calculated metrics based on these logs, select 5 from the menu (Read Metrics).
