import { fileURLToPath } from 'url';
import {} from 'dotenv/config';
import path from 'path';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export default function setConfig({ printConfig = false } = { printConfig: false }) {
	if (!process.env.ADMIN_API_KEY) {
		console.error('ENV ADMIN_API_KEY is required');
		throw new Error('ENV ADMIN_API_KEY is required');
	}

	if (!process.env.AUTH_URL) process.env.AUTH_URL = 'http://localhost:4000';

	if (!process.env.JWT_SHARED_SECRET) {
		console.error('ENV JWT_SHARED_SECRET is required');
		throw new Error('ENV JWT_SHARED_SECRET is required');
	}

	if (printConfig) {
		console.log('Starting with ENV:');
		console.log('ADMIN_API_KEY', '***')
		console.log('AUTH_URL', process.env.AUTH_URL);
		console.log('JWT_SHARED_SECRET', '***')
	}
}
