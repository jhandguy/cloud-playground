import http from 'k6/http';
import {check, sleep} from 'k6';
import {randomItem, randomString, uuidv4} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
    setupTimeout: '3m',
    scenarios: {
        load: {
            executor: 'ramping-arrival-rate',
            startRate: 1,
            timeUnit: '1s',
            preAllocatedVUs: 10,
            stages: [
                {target: 10, duration: '40s'},
                {target: 0, duration: '20s'},
            ],
        },
    },
    thresholds: {
        'checks': ['rate>0.9'],
        'http_req_duration{method:POST}': ['p(95)<10000'],
        'http_req_duration{method:GET}': ['p(95)<10000'],
        'http_req_duration{method:DELETE}': ['p(95)<10000'],
    },
};

const url = `http://${__ENV.SQL_URL}`;
const params = {
    headers: {
        'Authorization': `Bearer ${__ENV.SQL_TOKEN}`,
        'Host': __ENV.SQL_HOST,
        'Content-Type': 'application/json',
        'X-Redis-Enabled': __ENV.REDIS_ENABLED,
    },
};

export function setup() {
    let users = [];
    for (let i = 0; i < 3; i++) {
        let messages = [];
        const user = {
            id: uuidv4(),
            name: randomString(50),
        };

        check(http.post(`${url}/user`, JSON.stringify(user), params), {
            'post response status is not 403': (r) => r.status !== 403,
            'post response status is 201': (r) => r.status === 201,
            'post response body is valid': (r) => r.json().id === user.id && r.json().name === user.name,
        });

        for (let j = 0; j < 1000; j++) {
            const message = {
                id: uuidv4(),
                content: randomString(50),
                user_id: user.id,
            };

            check(http.post(`${url}/message`, JSON.stringify(message), params), {
                'post response status is not 403': (r) => r.status !== 403,
                'post response status is 201': (r) => r.status === 201,
                'post response body is valid': (r) => r.json().id === message.id && r.json().content === message.content && r.json().user_id === message.user_id,
            });

            messages.push(message)

            sleep(0.03);
        }

        users.push({
            id: user.id,
            name: user.name,
            messages: messages,
        });
    }

    return {
        users: users,
    }
}

export default function (data) {
    const user = randomItem(data.users);
    check(http.get(`${url}/user/${user.id}/messages`, params), {
        'get response status is not 403': (r) => r.status !== 403,
        'get response status is 200': (r) => r.status === 200,
        'get response body is valid': (r) => r.json().length === user.messages.length,
    });
}

export function teardown(data) {
    for (let user of data.users) {
        check(http.del(`${url}/user/${user.id}`, null, params), {
            'delete response status is not 403': (r) => r.status !== 403,
            'delete response status is 200': (r) => r.status === 200,
        });
    }
}
