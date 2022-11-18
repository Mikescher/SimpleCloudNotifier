
function send()
{
	let me = document.getElementById("btnSend");
	if (me.classList.contains("btn-disabled")) return;

	me.innerHTML = "<div class=\"spinnerbox\"><div class=\"spinner primary\"></div></div>";

	me.classList.add("btn-disabled");

	let uid = document.getElementById("uid");
	let key = document.getElementById("ukey");
	let msg = document.getElementById("msg");
	let txt = document.getElementById("txt");
	let pio = document.getElementById("prio");

	uid.classList.remove('input-invalid');
	key.classList.remove('input-invalid');
	msg.classList.remove('input-invalid');
	txt.classList.remove('input-invalid');
	pio.classList.remove('input-invalid');

	let data = new FormData();
	data.append('user_id', uid.value);
	data.append('user_key', key.value);
	data.append('title', msg.value);
	data.append('content', txt.value);
	data.append('priority', pio.value);

	let xhr = new XMLHttpRequest();
	xhr.open('POST', '/', true);
	xhr.onreadystatechange = function ()
	{
		if (xhr.readyState !== 4) return;

		console.log('Status: ' + xhr.status);
		if (xhr.status === 200 || xhr.status === 401 || xhr.status === 403 || xhr.status === 412)
		{
			let resp = JSON.parse(xhr.responseText);
			if (!resp.success || xhr.status !== 200)
			{
				if (resp.errhighlight === 101) uid.classList.add('input-invalid');
				if (resp.errhighlight === 102) key.classList.add('input-invalid');
				if (resp.errhighlight === 103) msg.classList.add('input-invalid');
				if (resp.errhighlight === 104) txt.classList.add('input-invalid');
				if (resp.errhighlight === 105) pio.classList.add('input-invalid');

				Toastify({
					text: resp.message,
					gravity: "top",
					positionLeft: false,
					backgroundColor: "#D32F2F",
				}).showToast();
			}
			else
			{
				window.location.href =
					'/message_sent' +
					'?ok=' + 1 +
					'&message_count=' + resp.messagecount +
					'&quota=' + resp.quota +
					'&quota_remain=' + (resp.quota_max-resp.quota) +
					'&quota_max=' + resp.quota_max +
					'&preset_user_id=' + uid.value +
					'&preset_user_key=' + key.value;
			}
		}
		else
		{
			Toastify({
				text: 'Request failed: Statuscode=' + xhr.status,
				gravity: "top",
				positionLeft: false,
				backgroundColor: "#D32F2F",
			}).showToast();
		}

		me.classList.remove("btn-disabled");
		me.innerHTML = "Send";
	};
	xhr.send(data);
}

window.addEventListener("load",function ()
{
	let btnSend = document.getElementById("btnSend");

	if (btnSend !== undefined) btnSend.onclick = function () { send(); return false; };

},false);