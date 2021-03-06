<!DOCTYPE html>
<html lang="en">
<?php
if (file_exists('/var/www/openwebanalytics/owa_php.php'))
{
	require_once('/var/www/openwebanalytics/owa_php.php');
	$owa = new owa_php();
	$owa->setSiteId('6386b0efc00d2e84ef642525345e1207');
	$owa->setPageTitle('Index');
	$owa->trackPageView();
}
?>
<head>
	<meta charset="utf-8">
    <title>Simple Cloud Notifications</title>
	<link rel="stylesheet" href="/css/toastify.min.css"/>
	<link rel="stylesheet" href="/css/mini-default.min.css">            <!-- https://minicss.org/ -->
	<!--<link rel="stylesheet" href="/css/mini-nord.min.css">-->
	<!--<link rel="stylesheet" href="/css/mini-dark.min.css">-->
	<link rel="stylesheet" href="/css/style.css">
	<meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" type="image/png" href="/favicon.png"/>
    <link rel="icon" type="image/png" href="/favicon.ico"/>
</head>
<body>

    <div id="copyinfo">
        <a tabindex="-1" href="https://www.blackforestbytes.com">&#169; blackforestbytes</a>
        <a tabindex="-1" href="https://www.mikescher.com">made by Mike Schw&ouml;rer</a>
    </div>

	<form id="mainpnl">

		<a tabindex="-1" href="https://play.google.com/store/apps/details?id=com.blackforestbytes.simplecloudnotifier" class="button bordered" id="tl_link"><span class="icn-google-play"></span></a>
		<a tabindex="-1" href="/api.php" class="button bordered" id="tr_link">API</a>

        <a tabindex="-1" href="/" class="linkcaption"><h1>Simple Cloud Notifier</h1></a>

		<div class="row responsive-label">
			<div class="col-sm-12 col-md-3"><label for="uid" class="doc">UserID</label></div>
			<div class="col-sm-12 col-md"><input placeholder="UserID" id="uid" class="doc" <?php echo (isset($_GET['preset_user_id']) ? (' value="'.$_GET['preset_user_id'].'" '):(''));?> type="number"></div>
		</div>

		<div class="row responsive-label">
			<div class="col-sm-12 col-md-3"><label for="ukey" class="doc">Authentification Key</label></div>
			<div class="col-sm-12 col-md"><input placeholder="Key" id="ukey" class="doc" <?php echo (isset($_GET['preset_user_key']) ? (' value="'.$_GET['preset_user_key'].'" '):(''));?> type="text" maxlength="64"></div>
		</div>

        <div class="row responsive-label">
            <div class="col-sm-12 col-md-3"><label for="prio" class="doc">Priority</label></div>
            <div class="col-sm-12 col-md">
                <select id="prio" class="doc" type="text" style="width:100%;">
                    <option value="0"  <?php echo (( isset($_GET['preset_priority'])&&$_GET['preset_priority']==='0') ? 'selected':'');?>>Low</option>
                    <option value="1"  <?php echo ((!isset($_GET['preset_priority'])||$_GET['preset_priority']==='1') ? 'selected':'');?>>Normal</option>
                    <option value="2"  <?php echo (( isset($_GET['preset_priority'])&&$_GET['preset_priority']==='2') ? 'selected':'');?>>High</option>
                </select>
            </div>
        </div>

		<div class="row responsive-label">
			<div class="col-sm-12 col-md-3"><label for="msg" class="doc">Message Title</label></div>
			<div class="col-sm-12 col-md"><input placeholder="Message" id="msg" class="doc" <?php echo (isset($_GET['preset_title']) ? (' value="'.$_GET['preset_title'].'" '):(''));?> type="text" maxlength="80"></div>
		</div>

		<div class="row responsive-label">
			<div class="col-sm-12 col-md-3"><label for="txt" class="doc">Message Content</label></div>
			<div class="col-sm-12 col-md"><textarea id="txt" class="doc" <?php echo (isset($_GET['preset_content']) ? (' value="'.$_GET['preset_content'].'" '):(''));?> rows="8" maxlength="2048"></textarea></div>
		</div>

		<div class="row">
			<div class="col-sm-12 col-md-3"></div>
			<div class="col-sm-12 col-md"><button type="submit" class="primary bordered" id="btnSend">Send</button></div>
		</div>
	</form>

	<script src="/js/logic.js" type="text/javascript" ></script>
	<script src="/js/toastify.js"></script>
</body>
</html>