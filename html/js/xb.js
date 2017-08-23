function LeonInit() {
	$('#editEvent').on('show.bs.modal', function (event) {
		var button = $(event.relatedTarget);
		var eventId = button.data('event-id');
		var modal = $(this);
		//ajax here
		if (0 != eventId) {
			$.ajax({
				url:'/msg/event-info',
				contentType:'application/json',
				data:JSON.stringify({Id:eventId}),
				type:'post',
				dataType:'json',
				success:function(data) {
					modal.find('#event-time').val(data.Time);
					modal.find('#event-desc').text(data.Desc);
				}
			});
		} else {
			modal.find('#event-time').val("");
			modal.find('#event-desc').text("");
		}
		modal.find('.modal-footer .btn-primary').attr('onclick', "saveEvent('"+eventId+"')");
	});

	var uploader = Qiniu.uploader({
		runtimes: 'html5',                  // 上传模式
		browse_button: $('#editEvent #event-image button'),         // 上传选择的点选按钮，必需
		uptoken_url: '/msg/qiniu-token',         // Ajax请求uptoken的Url，强烈建议设置（服务端提供）
		get_new_uptoken: false,             // 设置上传文件的时候是否每次都重新获取新的uptoken
		domain: 'xiaobai',     // bucket域名，下载资源时用到，必需
		container: $('#editEvent #event-image'),             // 上传区域DOM ID，默认是browser_button的父元素
		max_file_size: '100mb',             // 最大文件体积限制
		flash_swf_url: 'js/Moxie.swf',  //引入flash，相对路径
		max_retries: 3,                     // 上传失败最大重试次数
		dragdrop: false,                     // 关闭可拖曳上传
		chunk_size: '4mb',                  // 分块上传时，每块的体积
		auto_start: false,                  // 选择文件后自动上传，若关闭需要自己绑定事件触发上传
		init: {
			'FilesAdded': function(up, files) {
				plupload.each(files, function(file) {
					// 文件添加进队列后，处理相关的事情
				});
			},
			'BeforeUpload': function(up, file) {
				   // 每个文件上传前，处理相关的事情
			},
			'UploadProgress': function(up, file) {
				   // 每个文件上传时，处理相关的事情
			},
			'FileUploaded': function(up, file, info) {
				   // 每个文件上传成功后，处理相关的事情
				   // 其中info是文件上传成功后，服务端返回的json，形式如：
				   // {
				   //    "hash": "Fh8xVqod2MQ1mocfI4S4KpRL6D98",
				   //    "key": "gogopher.jpg"
				   //  }
				   // 查看简单反馈
				   var domain = up.getOption('domain');
				   var res = parseJSON(info);
				   var sourceLink = domain +"/"+ res.key; 获取上传成功后的文件的Url
			},
			'Error': function(up, err, errTip) {
				   //上传出错时，处理相关的事情
			},
			'UploadComplete': function() {
				   //队列文件处理完毕后，处理相关的事情
			},
			'Key': function(up, file) {
				// 若想在前端对每个文件的key进行个性化处理，可以配置该函数
				// 该配置必须要在unique_names: false，save_key: false时才生效

				var key = "";
				// do something with key here
				return key
			}
		}
	});
}
(function() {
	$(document).ready(function() {
		LeonInit();

		var timelineAnimate;
		timelineAnimate = function(elem) {
			return $(".timeline.animated .timeline-row").each(function(i) {
				var bottom_of_object, bottom_of_window;
				bottom_of_object = $(this).position().top + $(this).outerHeight();
				bottom_of_window = $(window).scrollTop() + $(window).height();
				if (bottom_of_window > bottom_of_object) {
					return $(this).addClass("active");
				}
			});
		};
		timelineAnimate();
		return $(window).scroll(function() {
			return timelineAnimate();
		});
	});
}).call(this);
