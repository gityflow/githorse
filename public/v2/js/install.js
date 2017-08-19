// Database type change detection.
$("#db_type").change(function () {
	var sqliteDefault = 'data/gogs.db';

	var dbType = $(this).val();
	if (dbType === "SQLite3") {
		$('#sql_settings').hide();
		$('#pgsql_settings').hide();
		$('#sqlite_settings').show();

		if (dbType === "SQLite3") {
			$('#db_path').val(sqliteDefault);
		}
		return;
	}

	var dbDefaults = {
		"MySQL": "127.0.0.1:3306",
		"PostgreSQL": "127.0.0.1:5432",
		"MSSQL": "127.0.0.1, 1433"
	};

	$('#sqlite_settings').hide();
	$('#sql_settings').show();
	$('#pgsql_settings').toggle(dbType === "PostgreSQL");
	$.each(dbDefaults, function (type, defaultHost) {
		if ($('#db_host').val() == defaultHost) {
			$('#db_host').val(dbDefaults[dbType]);
			return false;
		}
	});
});

// TODO: better handling of exclusive relations.
$('#offline-mode input').change(function () {
	if ($(this).is(':checked')) {
		$('#disable-gravatar').checkbox('check');
		$('#federated-avatar-lookup').checkbox('uncheck');
	}
});
$('#disable-gravatar input').change(function () {
	if ($(this).is(':checked')) {
		$('#federated-avatar-lookup').checkbox('uncheck');
	} else {
		$('#offline-mode').checkbox('uncheck');
	}
});
$('#federated-avatar-lookup input').change(function () {
	if ($(this).is(':checked')) {
		$('#disable-gravatar').checkbox('uncheck');
		$('#offline-mode').checkbox('uncheck');
	}
});
$('#disable-registration input').change(function () {
	if ($(this).is(':checked')) {
		$('#enable-captcha').checkbox('uncheck');
	}
});
$('#enable-captcha input').change(function () {
	if ($(this).is(':checked')) {
		$('#disable-registration').checkbox('uncheck');
	}
});

$('.ui.dropdown').dropdown();
$('.ui.checkbox').checkbox();
