import got from 'got';
import jwt from 'jsonwebtoken'
import setConfig from '../test-helpers/config.js';
import test from 'tape';

let adminJWT;
let adminJWTString;
let user;
let userJWT;
let userJWTString;
const userName = 'test-tomte nöff #18';
const password = 'lurpassare7½TUR';

test('test-cases/01basic.js: Authing with configurated API KEY', async t => {
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

	// Successful auth
	const authRes = await got.post(`${process.env.AUTH_URL}/auth/api-key`, {
		json: 'hihi',
		responseType: 'json',
	});
	t.notEqual(authRes.body.jwt, undefined, 'The body should include a jwt key');
	t.notEqual(authRes.body.renewalToken, undefined, 'The body should include a renewalToken');
	adminJWTString = authRes.body.jwt;

	adminJWT = jwt.verify(adminJWTString, process.env.JWT_SHARED_SECRET);
	t.equal(adminJWT.accountName, 'admin', 'The verified account name should be "admin"');
});

test('test-cases/01basic.js: GETting the admin account, with the token we just obtained', async t => {
	try {
		await got(`${process.env.AUTH_URL}/account/${adminJWT.accountId}`);
		t.fail('Calling /account/{id} without proper auth token should give 403');
	} catch (err) {
		t.equal(err.message, 'Response code 403 (Forbidden)', 'Calling /account/{id} without proper auth token should give 403');
	}

	const accountRes = await got(`${process.env.AUTH_URL}/account/${adminJWT.accountId}`, {
		headers: { 'Authorization': `bearer ${adminJWTString}`},
		responseType: 'json',
	});

	t.equal(adminJWT.accountId, accountRes.body.id, 'The account ids should match');
});

test('test-cases/01basic.js: Creating a new account', async t => {
	const res = await got.post(`${process.env.AUTH_URL}/account`, {
		headers: { 'Authorization': `bearer ${adminJWTString}`},
		json: {
			fields: [
				{
					name: 'nördområde',
					values: ['tåg', 'trädgårdstomtar'],
				},
				{
					name: 'role',
					values: ['user'],
				}
			],
			name: userName,
			password,
		},
		responseType: 'json',
	});

	user = res.body;

	t.notEqual(user.id, undefined, 'The new account should have an id');
	t.notEqual(user.apiKey, undefined, 'The new account should have an apiKey');

	try {
		await got.post(`${process.env.AUTH_URL}/account`, {
			headers: { 'Authorization': `bearer ${adminJWTString}`},
			json: {
				fields: [{name: 'role',values: ['user'],}],
				name: userName,
				password,
			},
			responseType: 'json',
		});
		t.fail('Trying to create another account with the same name should fail with a 409');
	} catch(err) {
		t.equal(err.message, 'Response code 409 (Conflict)', 'Trying to create another account with the same name should fail with a 409');
	}
});

test('test-cases/01basic.js: Auth by username and password', async t => {
	const authRes = await got.post(`${process.env.AUTH_URL}/auth/password`, {
		json: {
			name: userName,
			password,
		},
		responseType: 'json',
	});
	t.notEqual(authRes.body.jwt, undefined, 'The body should include a jwt key');
	t.notEqual(authRes.body.renewalToken, undefined, 'The body should include a renewalToken');
	userJWTString = authRes.body.jwt;

	userJWT = jwt.verify(userJWTString, process.env.JWT_SHARED_SECRET);
	t.equal(userJWT.accountName, userName, 'The verified account name should match the created user');
});

test('test-cases/01basic.js: Remove an account', async t => {
	try {
		// Random uuid that should not exist in the db. The chance of this existing is... small
		await got.delete(`${process.env.AUTH_URL}/account/a423e690-74b9-4f37-9976-f5bf75a5ea32`, {
			headers: { 'Authorization': `bearer ${adminJWTString}`},
			responseType: 'json',
			retry: 0,
		});
		t.fail('Response status for DELETing an account that does not exist should be 404');
	} catch (err) {
		t.equal(err.message, 'Response code 404 (Not Found)', 'Response status for DELETing an account that does not exist should be 404');
	}

	const delRes = await got.delete(`${process.env.AUTH_URL}/account/${user.id}`, {
		headers: { 'Authorization': `bearer ${adminJWTString}`},
		responseType: 'json',
		retry: 0,
	});

	t.equal(delRes.statusCode, 204, 'Response status for DELETE should be 204');

	try {
		await got(`${process.env.AUTH_URL}/account/${user.id}`, {
			headers: { 'Authorization': `bearer ${adminJWTString}`},
			responseType: 'json',
			retry: 0,
		});
		t.fail('Response status for GETing the account again should be 404');
	} catch (err) {
		t.equal(err.message, 'Response code 404 (Not Found)', 'Response status for GETing the account again should be 404');
	}
});
