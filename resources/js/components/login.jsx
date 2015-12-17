var Comment = React.createClass({
  render: function() {
    return (
      <div className="comment">
        <h2 className="commentAuthor">{this.props.author}</h2>
        {this.props.children}
      </div>
    );
  }
});

var CommentList = React.createClass({

  render: function() {
    var commentNodes = this.props.data.map(function(comment) {
      return (
        <Comment author={comment.author} key={comment.id}>
          {comment.text}
        </Comment>
      );
    });

    return (
      <div className="commentList">{commentNodes}</div>
    );
  }
});

var CommentBox = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  componentDidMount: function() {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
    })
    .done(function(data) {
      this.setState({data: data});
    }.bind(this))
    .fail(function(xhr, status, err) {
      
    }.bind(this));

  },
  render: function() {
    return (
      <div className="commentBox">
        <h1>Comments</h1>
        <CommentList data={this.state.data}/>
        <CommentForm/>
      </div>
    );
  }
});

var CommentForm = React.createClass({

  render: function() {
    return (
      <div className="commentForm">
        Hello, world! I am a CommentForm.
      </div>
    );
  }
});

ReactDOM.render(
  <CommentBox url="/demos" />,
  document.getElementById('app-container'));

// var LoginForm = React.createClass({
//   render: function() {
//     return (
//       <form className="form-signin" method="post">
//         <h2 className="form-signin-heading">Please sign in</h2>

//         <label htmlFor="inputEmail" className="sr-only">Email address</label>
//         <input type="email" id="inputEmail" className="form-control" placeholder="Email address" required autofocus/>

//         <label htmlFor="inputPassword" className="sr-only">Password</label>
//         <input type="password" id="inputPassword" className="form-control" placeholder="Password" required/>

//         <button className="btn btn-lg btn-primary btn-block" type="submit" onClick="">Sign in</button>
//       </form>
//     );
//   }
// });
// ReactDOM.render(< LoginForm />, document.getElementById('container'));
