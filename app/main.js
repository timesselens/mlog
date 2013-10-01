/** @jsx React.DOM */
var Console = React.createClass({
    getInitialState: function() {
        return {data: this.props.data};
    },
    render: function() {
        var html = this.state.data.map(function(o,i) {
            return <div class="logline" dangerouslySetInnerHTML={{__html: o.html}}></div>;
        });

        return <div class="console">{html}</div>

    }
});

window.data = []; // JSON.parse(localStorage['last10'])

var consoleview = Console({data: window.data});
//React.renderComponent(consoleview, $('.console-container')[0]);

var ws = new WebSocket("ws://localhost:12345/")
ws.onmessage = function(m) {
    // console.log('message %s',m.data)
    var tpl = $('<div class="logline"></div>').text(m.data)
    $('div.console').append(tpl);
    $('div.console-container').scrollTop(99999);
}

localStorage['mlog_config'] = JSON.stringify({
    procs: ["stdin://", { type: "grep", re: /^[\w\d]+$/.toString() }], // accepts both strings and hashes
    pipes: [[0,1],
            [1,2]]
});

$.post('/setup?id=main', localStorage['mlog_config'])
