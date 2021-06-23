import got from 'got';
import jwt from 'jsonwebtoken'
import setConfig from '../test-helpers/config.js';
import test from 'tape';

test('test-cases/01basic.js: Basic stuff', async t => {
	t.comment('Authing with configurated API KEY');

	// Wrong API key
	try {
		await got.post(`${process.env.AUTH_URL}/auth/api-key`, {
			json: 'a09ifa908wjf92fowreigaoijfaosidfđ@€£đawef',
			responseType: 'json',
		});

		t.fail('Calling /auth/api-key with wrong api-key should result in a 403');
	} catch (err) {
		t.equal(err.message, 'Response code 403 (Forbidden)', 'Calling /auth/api-key with wrong api-key should result in a 403')
	}

	const authRes = await got.post(`${process.env.AUTH_URL}/auth/api-key`, {
		json: 'hihi',
		responseType: 'json',
	});
	t.notEqual(authRes.body.jwt, undefined, 'The body should include a jwt key');
	t.notEqual(authRes.body.renewalToken, undefined, 'The body should include a renewalToken');

	const adminJWT = jwt.verify(authRes.body.jwt, process.env.JWT_SHARED_SECRET);
	t.equal(adminJWT.accountName, 'admin', 'The verified account name should be "admin"');

	t.comment('GETting the admin account, with the token we just obtained');

	try {
		await got(`${process.env.AUTH_URL}/account/${adminJWT.accountId}`);
		t.fail('Calling /account/{id} without proper auth token should give 403');
	} catch (err) {
		t.equal(err.message, 'Response code 403 (Forbidden)', 'Calling /account/{id} without proper auth token should give 403');
	}

	const accountRes = await got(`${process.env.AUTH_URL}/account/${adminJWT.accountId}`, {
		headers: { 'Authorization': `bearer ${authRes.body.jwt}`},
		responseType: 'json',
	});

	t.equal(adminJWT.accountId, accountRes.body.id, 'The account ids should match');

	t.end();
});
