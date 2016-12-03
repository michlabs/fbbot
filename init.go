package fbbot

const (
	WebhookURL      = "/webhook"
	SendAPIEndpoint = "https://graph.facebook.com/v2.6/me/messages"
	APIEndpoint     = "https://graph.facebook.com/v2.6"

	// Notification type
	NotiRegular    string = "REGULAR"     // will emit a sound/vibration and a phone notification
	NotiSilentPush string = "SILENT_PUSH" // will just emit a phone notification
	NotiNoPush     string = "NO_PUSH"     // will not emit either
)

const TEMPLATE = `<!doctype html>

<meta charset="utf-8">
<title>FBbot - Facebook Messenger bot framework for Humans, for Go</title>

<script src="http://d3js.org/d3.v3.min.js" charset="utf-8"></script>
<script src="http://cpettitt.github.io/project/dagre-d3/latest/dagre-d3.min.js"></script>

<style id="css">
body {
  font: 300 14px 'Helvetica Neue', Helvetica;
}

.node rect {
  stroke: #333;
  fill: #fff;
}

.edgePath path {
  stroke: #333;
  fill: #333;
  stroke-width: 1.5px;
}
</style>

<svg width=960 height=600><g/></svg>

<section>
<p>Dialog Diagram</p>
</section>

<script id="js">
// Create a new directed graph
var g = new dagreD3.graphlib.Graph().setGraph({});

// Set nodes 
%s

// Automatically label each of the nodes
states.forEach(function(state) { g.setNode(state, { label: state }); });

// Set Edges 
%s

// Set some general styles
g.nodes().forEach(function(v) {
  var node = g.node(v);
  node.rx = node.ry = 5;
});

// Add some custom colors based on state
g.node('%s').style = "fill: #7f7";
g.node('%s').style = "fill: #f77";

var svg = d3.select("svg"),
    inner = svg.select("g");

// Set up zoom support
var zoom = d3.behavior.zoom().on("zoom", function() {
      inner.attr("transform", "translate(" + d3.event.translate + ")" +
                                  "scale(" + d3.event.scale + ")");
    });
svg.call(zoom);

// Create the renderer
var render = new dagreD3.render();

// Run the renderer. This is what draws the final graph.
render(inner, g);

// Center the graph
var initialScale = 0.75;
zoom
  .translate([(svg.attr("width") - g.graph().width * initialScale) / 2, 20])
  .scale(initialScale)
  .event(svg);
svg.attr('height', g.graph().height * initialScale + 40);
</script>`