import grpc from 'k6/net/grpc';
import {check, sleep} from 'k6';
import {randomString, uuidv4} from 'https://jslib.k6.io/k6-utils/1.1.0/index.js';

export const options = {
    stages: [
        {target: 10, duration: '20s'},
        {target: 10, duration: '20s'},
        {target: 0, duration: '20s'},
    ],
    thresholds: {
        'checks': ['rate>0.9'],
        'grpc_req_duration{method:CreateObject}': ['p(95)<1000'],
        'grpc_req_duration{method:GetObject}': ['p(95)<1000'],
        'grpc_req_duration{method:DeleteObject}': ['p(95)<1000'],
    },
};

const client = new grpc.Client();
client.load(['../proto/object'], 'object.pg.proto');

export default function () {
    const params = {
        headers: {
            'x-token': __ENV.S3_TOKEN,
        },
    };
    const id = uuidv4()
    const content = randomString(50)
    const data = {
        object: {
            id: id,
            content: content,
        }
    };

    client.connect(__ENV.S3_URL, {
        plaintext: true,
    });

    check(client.invoke('object.ObjectService/CreateObject', data, params), {
        'create status is OK': (r) => r.status === grpc.StatusOK,
        'create message is valid': (r) => r.message.object.id === id && r.message.object.content === content,
    });

    sleep(0.5)

    check(client.invoke('object.ObjectService/GetObject', {id: id}, params), {
        'get status is OK': (r) => r.status === grpc.StatusOK,
        'get message is valid': (r) => r.message.object.id === id && r.message.object.content === content,
    });

    sleep(0.5)

    check(client.invoke('object.ObjectService/DeleteObject', {id: id}, params), {
        'delete status is OK': (r) => r.status === grpc.StatusOK,
    });

    sleep(1)
}
