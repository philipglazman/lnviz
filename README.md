
# Lightning History Visualization

##### DRAFT.

Visualize your routing node's performance on the Lightning Network!

Connect to your LND node, and create a page of reports to help you make better decisions.

## Usage

`go run . -cert="<filepath>" -macaroon="<filepath>" -host=<host:port>`

`open report.html`

## Examples
todo

## Remaining work
* Expose report as a server and subscribe to routing events to build a live reporting tool.
* Report: Interactive graph to see which the capacity between remote nodes.