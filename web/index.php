<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="/css/toastify.min.css"/>
	<link rel="stylesheet" href="/css/mini-default.min.css">
	<!--<link rel="stylesheet" href="/css/mini-nord.min.css">-->
	<!--<link rel="stylesheet" href="/css/mini-dark.min.css">-->
	<link rel="stylesheet" href="/css/style.css">
	<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
	<form id="mainpnl">

		<a href="https://play.google.com/store/apps/details?id=com.blackforestbytes.simplecloudnotifier" class="button bordered" id="tl_link"><span class="icn-google-play"></span></a>
		<a href="/index_api.php" class="button bordered" id="tr_link">API</a>

		<h1>Simple Cloud Notifier</h1>

		<div class="row responsive-label">
			<div class="col-sm-12 col-md-3"><label for="uid" class="doc">UserID</label></div>
			<div class="col-sm-12 col-md"><input placeholder="UserID" id="uid" class="doc" <?php echo (isset($_GET['preset_user_id']) ? (' value="'.$_GET['preset_user_id'].'" '):(''));?> type="number"></div>
		</div>

		<div class="row responsive-label">
			<div class="col-sm-12 col-md-3"><label for="ukey" class="doc">Authentification Key</label></div>
			<div class="col-sm-12 col-md"><input placeholder="Key" id="ukey" class="doc" <?php echo (isset($_GET['preset_user_key']) ? (' value="'.$_GET['preset_user_key'].'" '):(''));?> type="text" maxlength="64"></div>
		</div>

		<div class="row responsive-label">
			<div class="col-sm-12 col-md-3"><label for="msg" class="doc">Message Title</label></div>
			<div class="col-sm-12 col-md"><input placeholder="Message" id="msg" class="doc" <?php echo (isset($_GET['preset_title']) ? (' value="'.$_GET['preset_title'].'" '):(''));?> type="text" maxlength="80"></div>
		</div>

		<div class="row responsive-label">
			<div class="col-sm-12 col-md-3"><label for="txt" class="doc">Message Content</label></div>
			<div class="col-sm-12 col-md"><textarea id="txt" class="doc" <?php echo (isset($_GET['preset_content']) ? (' value="'.$_GET['preset_content'].'" '):(''));?> rows="5"></textarea></div>
		</div>

		<div class="row">
			<div class="col-sm-12 col-md-3"></div>
			<div class="col-sm-12 col-md"><button type="submit" class="primary bordered" id="btnSend">Send</button></div>
		</div>
	</form>

	<div id="copyinfo">
		<a href="https://www.blackforestbytes.com">&#169; blackforestbytes</a>
	</div>

	<script src="/js/logic.js" type="text/javascript" ></script>
	<script src="/js/toastify.js"></script>
</body>
</html>