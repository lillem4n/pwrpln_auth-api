import {} from 'dotenv/config';
import got from 'got';
import test from 'tape';
import setConfig from '../test-helpers/config.js';

test('test-cases/00start.js: Wait for auth API to be ready', async t => {
	setConfig({ printConfig: true });

	const backendHealthCheck = await got(process.env.AUTH_URL, { retry: 2000 });

	t.equal(backendHealthCheck.statusCode, 200, 'Auth API should answer with status code 200');
});
