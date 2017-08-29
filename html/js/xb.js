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
				async: false,
				success:function(data) {
					modal.find('#event-time').val(data.Time);
					modal.find('#event-address').val(data.Address);
					modal.find('#event-title').val(data.Title);
					modal.find('#event-image img').attr('src', data.Image);
					modal.find('#event-desc').val(data.Desc);
					haveChgd = new Array();
				}
			});
		} else {
			modal.find('#event-time').val("");
			modal.find('#event-address').val("");
			modal.find('#event-title').val("");
			modal.find('#event-image img').attr('src', '');
			modal.find('#event-desc').val("");
			haveChgd = new Array();
		}
		// if (modal.find('.modal-header #xb-event-id').text() != eventId) {
		// 	modal.find('.modal-header #xb-event-id').text(eventId);
		// 	modal.find('#chg-img').text(0);
		// }
		modal.find('.modal-header #xb-event-id').text(eventId);
		modal.find('.modal-footer .btn-primary').attr('onclick', "putSave()");
	});

	uploader = Qiniu.uploader({
		runtimes: 'html5',
		browse_button: 'event-image-add',
		// uptoken_url: '/msg/qiniu-token',  // Ajax请求uptoken的Url，格式过于死板
		uptoken_func: function() {
			// var token = "";
			$.ajax({
				url:'/msg/qiniu-token',
				contentType: 'application/json',
				data:JSON.stringify({Bucket:'xiaobai'}),
				type:'post',
				dataType:'json',
				async:false,
				success:function(data) {
					token = data.Token;
				}
			})
			return token;
		},
		unique_names: true,
		get_new_uptoken: false,
		domain: 'xiaobai',
		container: 'event-image',
		max_file_size: '100mb',
		max_retries: 3,
		chunk_size: '4mb',
		multi_selection: false,
		filters: {
			mime_types: [
				{title:"图片文件", extensions:"jpg,png"}
			]
		},
		auto_start: false,
		init: {
			'FilesAdded': function(up, files) {
				plupload.each(files, function(file) {
					// 文件添加进队列后，处理相关的事情
					console.log("add file:", file.name);
					var preloader = new mOxie.Image();
					preloader.onload = function() {
						preloader.downsize(300, 300); //压缩下显示 不影响上传
						var imgsrc = preloader.type=='image/jpeg' ? preloader.getAsDataURL('image/jpeg',80) : preloader.getAsDataURL();
						$('#editEvent #event-image img').attr('src', imgsrc);
						preloader.destroy();
						preloader = null;
					};
					preloader.load(file.getSource());
					// $('#editEvent #event-image #chg-img').text(1);
					if (-1 === haveChgd.indexOf('event-image'))
						haveChgd.push('event-image');
				});
			},
			'FileUploaded': function(up, file, info) {
				var domain = up.getOption('domain');
				var res = JSON.parse(info.response);
				var sourceLink = domain +"/"+ res.key; //上传成功后 域+名
				$('#editEvent #event-image img').attr('src', sourceLink);
				saveEvent($('#editEvent .modal-header #xb-event-id').text());
			},
			'Error': function(up, err, errTip) {
				//上传出错时，处理相关的事情
			}
		}
	});

	function bindChg() {
		var chgId = $(this).attr('id');
		if (-1 === haveChgd.indexOf(chgId))
			haveChgd.push(chgId);
	}
	$('#editEvent #event-time').bind('input propertychange', bindChg);
	$('#editEvent #event-address').bind('input propertychange', bindChg);
	$('#editEvent #event-title').bind('input propertychange', bindChg);
	$('#editEvent #event-desc').bind('input propertychange', bindChg);

	laydate.render({
		elem: '#editEvent #event-time'
	});
}
function putSave() {
	// if (1 == $('#editEvent #event-image #chg-img').text()) {
	if (-1 === haveChgd.indexOf('event-image')) {
		saveEvent($('#editEvent .modal-header #xb-event-id').text());
	} else {
		uploader.start();
	}
}
var haveChgd = new Array();
function saveEvent(eventId) {
	var data = {Id:eventId}//new Object();
	var dom = haveChgd.pop();
	while (undefined !== dom) {
		switch (dom) {
			case 'event-time': {
				data.Time = $('#editEvent #event-time').val();
				break;
			}
			case 'event-address': {
				data.Address = $('#editEvent #event-address').val();
				break;
			}
			case 'event-title': {
				data.Title = $('#editEvent #event-title').val();
				break;
			}
			case 'event-image': {
				data.Image = $('#editEvent #event-image img').attr('src');
				break;
			}
			case 'event-desc': {
				data.Desc = $('#editEvent #event-desc').val();
				break;
			}
			default: {
				break;
			}
		}
		dom = haveChgd.pop();
	}
	// console.log(data);
	$.ajax({
		url: '/msg/event-save',
		contentType: 'application/json',
		data: JSON.stringify(data),
		type: 'post',
		dataType: 'json',
		success:function(data) {
			if (data.Result)
				window.location.reload();
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
