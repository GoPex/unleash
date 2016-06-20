var PingDockerHost = React.createClass({
  ping: function() {
    this.setState({status: "Pinging ..."});
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      type: 'GET',
      cache: false,
      success: function(data) {
        this.setState({status: data.pong});
      }.bind(this),
      error: function(xhr, status, err) {
        this.setState({status: "NOK"});
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },

  handleClick: function(e) {
    e.preventDefault();
    this.ping()
  },

  getInitialState: function() {
    return {status: ""};
  },

  componentDidMount: function() {
    this.ping();
  },

  render: function() {
    return (
      <a className="pingDockerHost" onClick={this.handleClick} href={this.props.url}>Ping your Docker host ! Status : {this.state.status}</a>
    );
  }
});

ReactDOM.render(
  <PingDockerHost url="/ping" />,
  document.getElementById("ping-docker-host")
);
