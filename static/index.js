var request = window.superagent;

new Vue({
	el: "#contents",
	data: {
		posts: [],
		formContents: ""
	},
	methods: {
		createPost: function (contents) {
			if(contents.length === 0){
				return false;
			}
			
			var self = this;
			request.post("/api/posts/")
				.type('form')
				.send({ Contents: contents })
				.end(function (err, res) {
					self.readPosts();
					self.formContents = "";
				});
		},
		readPosts: function () {
			var self = this;
			request.get("/api/posts/")
				.end(function (err, res) {
					var result = res.body;
					result.reverse();
					self.posts = result;
				});
		}
	},
	ready: function () {
		this.readPosts();
	}
})