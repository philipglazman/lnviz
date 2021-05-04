
# lnviz: Lightning Network Visualization

##### DRAFT.

Visualize your routing node's performance on the Lightning Network.

Connect to your LND node, and generate a report to help you make better decisions on channel management, liquidity, and peers.

## Usage
lnviz creates a file titled `report.html` containing a report with charts.

`go run . -cert="<filepath>" -macaroon="<filepath>" -host=<host:port>`

`open report.html`

## Examples


## Remaining work
* Expose report as a server and subscribe to routing events to build a live reporting tool.
* Report: Interactive graph to see which the capacity between remote nodes.