
# lnviz: Lightning Network Visualization

##### DRAFT.

Visualize your routing node's performance on the Lightning Network.

Connect to your LND node, and generate a report to help you make better decisions on channel management, liquidity, and peers.

## Usage
lnviz creates a file titled `report.html` containing a report with charts.

`go run . -cert="<filepath>" -macaroon="<filepath>" -host=<host:port>`

`open report.html`

## Examples
![sum_route_fee_outbound](https://user-images.githubusercontent.com/8378656/116965885-c910e700-ac63-11eb-9696-3aa2fe0ddd27.png)

![sum_route_fee](https://user-images.githubusercontent.com/8378656/116965895-d1692200-ac63-11eb-99cc-7f4365b55452.png)

![sum_rout_fee_inbound](https://user-images.githubusercontent.com/8378656/116965901-d4fca900-ac63-11eb-8df9-2f3c7f5c5fb9.png)

## Remaining work
* Expose report as a server and subscribe to routing events to build a live reporting tool.
* Report: Interactive graph to see which the capacity between remote nodes.
