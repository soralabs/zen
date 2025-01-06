package random_toolkit

import toolkit "github.com/soralabs/toolkit/go"

var Toolkit = toolkit.NewToolkit("random_tools",
	toolkit.WithToolkitDescription("A toolkit that include random generation"),
	toolkit.WithTools(
		NewRandomNumberTool(),
	),
)
