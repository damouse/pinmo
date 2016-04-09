package com.ironether.client;

import java.io.BufferedReader;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.InetAddress;
import java.net.Socket;

/**
 * This is the "client"
 * it should call into Go code in a way thats:
 *  - type safe
 *  - "static" wrt function calls, by dynamic wrt the exposed functionality
 *  - asynchronous
 */
public class Main {
    public static void log(String s) {
        System.out.println(s);
    }

    public static void main(String[] args) {
        log("Starting the client");

        MrsFrizzle f = new MrsFrizzle("localhost", 9876);
        f.send("Java: Hi!\n");

        f.run();
    }
}

// Its a magic bus. Get it?
class MrsFrizzle {
    DataInputStream in;
    DataOutputStream out;
    Thread thread;

    MrsFrizzle(String host, int port) {
        Main.log("Connecting to " + host + " " + port);
        try {
            Socket socket = new Socket(InetAddress.getByName(host), port);
            in = new DataInputStream(socket.getInputStream());
            out = new DataOutputStream(socket.getOutputStream());
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    // Spin and listen for events
    public void listen() {
        thread = new Thread() {
            public void run() {
                    run();
                }

        };
        thread.start();
    }

    // Spin and listen
    public void run() {
        BufferedReader d = new BufferedReader(new InputStreamReader(in));

        while (true) {
            try {
                String line = d.readLine();

                if (line == null)
                    break;

                System.out.println("Java has: " + line);

            } catch (Exception e) {
                Main.log("Things didnt go well");
                e.printStackTrace();
            }
        }

        Main.log("Connection lost");
    }

    public void send(String response) {
        Main.log("Writing " + response);

        try {
            out.writeBytes(response);
            out.flush();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}