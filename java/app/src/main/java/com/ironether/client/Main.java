package com.ironether.client;

import com.google.gson.Gson;

import java.io.BufferedReader;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.math.BigInteger;
import java.net.InetAddress;
import java.net.Socket;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Random;

/**
 * This is the "client"
 * it should call into Go code in a way thats:
 *  - type safe
 *  - "static" wrt function calls, by dynamic wrt the exposed functionality
 *  - asynchronous
 *
 *  The core only calls us with Deferred style callbacks- (uint64, uint64, args)
 */
public class Main {
    public static void log(String s) {
        System.out.println(s);
    }

    public static void main(String[] args) {
        log("Starting the client");

        MrsFrizzle f = new MrsFrizzle("localhost", 9876);

        f.affect("Info", "This is Patrick");

        try {
            f.thread.join();
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }
}

// Its a magic bus. Get it?
class MrsFrizzle {
    DataInputStream in;
    DataOutputStream out;
    Thread thread;

    Gson gson = new Gson();
    Map<Integer, Deferred> deferreds =  new HashMap();


    MrsFrizzle(String host, int port) {
        Main.log("Connecting to " + host + " " + port);

        try {
            Socket socket = new Socket(InetAddress.getByName(host), port);
            in = new DataInputStream(socket.getInputStream());
            out = new DataOutputStream(socket.getOutputStream());
        } catch (IOException e) {
            e.printStackTrace();
        }

        listen();
    }

    // Spin and listen for events
    public void listen() {
        thread = new Thread() {
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

        };
        thread.start();
    }


    public void send(String response) {
        Main.log("Writing " + response);

        try {
            out.writeBytes(response + "\n");
            out.flush();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    // Affect some change in the core. The result of the operation depends on what the type of the thing is.
    // Types can be instantiated and return a pointer
    // Functions can be invoked
    // Variables can be set and read
    // Constants can be read
    // Always returns a Deferred. The deferred is only invoked if its callback methods are called
    public Deferred affect(String target, Object... args) {
        List<Object> r = new ArrayList();
        Deferred d = new Deferred();
        r.add(target);
        r.add(d.cb);
        r.add(d.eb);
        r.add(args);
        send(gson.toJson(r));

        return d;
    }
}

class Deferred {
    int cb;
    int eb;
    Function _callback;
    Function _errback;

    private static Random gen = new Random();

    public Deferred() {
        cb = gen.nextInt();
        eb = gen.nextInt();
    }


    // We will need to override cuminicable methods here again, much like swift
    Deferred _then(Function fn) {
        _callback = fn;
        return this;
    }

    Deferred _error(Function fn) {
        _errback = fn;
        return this;
    }

    void callback(Object[] args) {
        if (_callback != null) {
            _callback.invoke(args);
        }
    }

    void errback(Object[] args) {
        if (_errback != null) {
            _errback.invoke(args);
        }
    }

    //
    // Generic Shotgun
    //

    // No args
    public Deferred then(Function handler) {
        return _then(handler);
    }

    public Deferred error(Function handler) {
        return _error(handler);
    }
}

interface Function {
    void invoke(Object[] args);
}