<?php

function prepareResponseSocket($responseSocketPath)
{
    $responseSocket = socket_create(AF_UNIX, SOCK_STREAM, 0);
    socket_bind($responseSocket, $responseSocketPath);
    socket_listen($responseSocket);
    return $responseSocket;
}

function sendToAgent(string $responseSocketPath, $messageContent)
{
    $agentSocketPath = '/tmp/rabbitmqagent.sock';
    $agentSocket = socket_create(AF_UNIX, SOCK_STREAM, 0);

    socket_connect($agentSocket, $agentSocketPath);
    socket_write($agentSocket, $messageContent, strlen($messageContent));
    socket_close($agentSocket);
}

function receiveFromAgent($responseSocket): string
{

    $sockets = [$responseSocket];
    $writeSockets = null;
    $except = null;

    socket_select($sockets, $writeSockets, $except, 20);

    $client = socket_accept($responseSocket);
    socket_set_nonblock($client);
    $buffer = socket_read($client, 2048);

    socket_close($client);
    socket_close($responseSocket);

    return $buffer;
}

$responseSocketPath = '/tmp/response_'.uniqid().'.sock';
$responseSocket = prepareResponseSocket($responseSocketPath);

$data = [
        'payload' => 'message payload',
        'deliver_options' => [
            'vhost' => 'okahkdhy',
            'exchange_name' => 'exchangename',
            'routing_key' => 'routing.key'
        ],
        'response_socket' => $responseSocketPath
    ];



sendToAgent($responseSocketPath, json_encode($data));
$response = receiveFromAgent($responseSocket);
unlink($responseSocketPath);

echo "Received response: $response" . PHP_EOL;





