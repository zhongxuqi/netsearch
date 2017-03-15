angular.module('app',['ui.router'])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state("default", {
				url: "/default",
				templateUrl: "/html/default.html",
				controller: "defaultControler",
			})
			.state("list_result", {
				url: "/list_result/:keyword/:ids",
				templateUrl: "/html/list_result.html",
				controller: "listControler",
			});
		$urlRouterProvider.when("", "/default");
	});

(function() {
	var totalSearchers = [ {
			id: 0,
			name: "百度",
			parser_name: "BaiduData",
			state: "active",
			result_list: [],
			next_page: ""
		}, {
			id: 1,
			name: "知乎",
			parser_name: "ZhihuData",
			state: "",
			result_list: [],
			next_page: ""
		}, {
			id: 2,
			name: "好搜",
			parser_name: "HaosouData",
			state: "",
			result_list: [],
			next_page: ""
		}, {
			id: 3,
			name: "维基百科",
			parser_name: "WikipediaData",
			state: "",
			result_list: [],
			next_page: ""
		}, {
			id: 4,
			name: "百度学术",
			parser_name: "BaiduXueShuData",
			state: "",
			result_list: [],
			next_page: ""
		}, {
			id: 5,
			name: "简书",
			parser_name: "JianShuData",
			state: "",
			result_list: [],
			next_page: ""
		}, {
			id: 6,
			name: "CSDN",
			parser_name: "CSDNData",
			state: "",
			result_list: [],
			next_page: ""
		}
	];
	angular.module('app').service("spiderServer", ["$http", spiderServer]);
	function spiderServer($http) {
		function getData(spiderId, keyword, callback, error) {
			var url = "";
			switch (spiderId) {
				case 0:
					url = "../app/search_baidu?keyword=";
					break;
				case 1:
					url = "../app/search_zhihu?keyword=";
					break;
				case 2:
					url = "../app/search_haosou?keyword=";
					break;
				case 3:
					url = "../app/search_wikipedia?keyword=";
					break;
				case 4:
					url = "../app/search_baiduxueshu?keyword=";
					break;
				case 5:
					url = "../app/search_jianshu?keyword=";
					break;
				case 6:
					url = "../app/search_csdn?keyword=";
					break;
			}
			$http({
				id: spiderId,
				url:url + keyword,
				method:'GET'
			}).success(callback).error(error);
		}
		function getNextPage(url, parser_name, callback, error) {
			$http({
				url: "../app/custom_search?parser_name="+parser_name+"&url="+url,
				method:'GET'
			}).success(callback).error(error);
		}
		return {
			getData: getData,
			getNextPage: getNextPage,
		};
	}
	angular.module('app').controller("defaultControler", ["$scope", "$state", defaultControler]);
	function defaultControler($scope, $state) {
		$scope.searcherList = [];
		for (var i = 0; i < totalSearchers.length; i++) {
			var item = {
				id: totalSearchers[i].id,
				name: totalSearchers[i].name,
				state: ""
			}
			if (i == 0 || i == 1|| i ==2) {
				item.state="active";
			}
			$scope.searcherList.push(item);
		}

		$scope.search = function(keyword) {
			if (keyword == undefined || keyword.length == 0) return;
			var ids = [];
			for (var i = 0; i < $scope.searcherList.length; i++) {
				if ($scope.searcherList[i].state == "active") {
					ids.push($scope.searcherList[i].id);
				}
			}
			$state.go("list_result", {"keyword": keyword, "ids": ids});
		}
		$scope.onTagClick = function(item) {
			if (item.state === "active") {
				item.state = "";
			} else {
				item.state = "active";
			}
		}
	}
	angular.module('app').controller("listControler", ["$scope", "$stateParams", "spiderServer", listControler]);
	function listControler($scope, $stateParams, spiderServer) {
		$scope.Searchers = [];
		var ids = $stateParams.ids.split(",");
		var cnt = 0;
		for (var i = 0; i < totalSearchers.length; i++) {
			for (var j = 0; j < ids.length; j ++) {
				if (totalSearchers[i].id == parseInt(ids[j])) {
					$scope.Searchers.push(totalSearchers[i]);
					$scope.Searchers[cnt].state = "";
					cnt++;
					break;
				}
			}
		}
		if (cnt > 0) {
			$scope.Searchers[0].state = "active";
		}
		$scope.switchState = function(Searcher) {
			for (var i = 0; i < $scope.Searchers.length; i++) {
				if (Searcher.id == $scope.Searchers[i].id) {
					$scope.currIndex = i;
					break;
				}
			}
			for (var i = 0; i < $scope.Searchers.length; i++) {
				if (Searcher.id === $scope.Searchers[i].id) {
					$scope.Searchers[i].state = "active";
				} else {
					$scope.Searchers[i].state = "";
				}
			}
		};
		$scope.search = function(keyword) {
			if (keyword == undefined || keyword.length == 0) return;
			for (var i = 0; i < $scope.Searchers.length; i++) {
				var item = $scope.Searchers[i];
				spiderServer.getData(item.id, keyword, function(response, status, headers, config) {
					for (var i = 0; i < $scope.Searchers.length; i++) {
						if (config.id == $scope.Searchers[i].id) {
							$scope.Searchers[i].result_list = response.ResultData;
							$scope.Searchers[i].next_page = response.NextPage;
							break;
						}
					}
				}, function(response, status, headers, config) {
					console.log(response);
				});
			}
		}
		$('#zxq-loading-more').on('click', function () {
			var currItem = $scope.Searchers[$scope.currIndex];
			var $btn = $(this).button('loading');
			if (currItem.next_page == undefined || currItem.next_page.length == 0) return;
			spiderServer.getNextPage(currItem.next_page, currItem.parser_name, function(response, status, headers, config) {
				try {
					for (var i = 0; i < response.ResultData.length; i++) {
						currItem.result_list.push(response.ResultData[i]);
					}
				} catch(e) {

				}
				currItem.next_page = response.NextPage;
				$btn.button('reset');
			}, function(response, status, headers, config) {
				console.log(response);
				$btn.button('reset');
			}) ;
		});
		if ($stateParams.keyword !== undefined && $stateParams.keyword.length > 0) {
			$scope.keyword = $stateParams.keyword;
			$scope.search($stateParams.keyword);
		}
		$scope.currIndex = 0; 
	}
})();
