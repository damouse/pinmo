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


        // Goal... but first a messaging bus
        // Go.call("Bark", 1).then -> result

        // We have inverted control here, so establish an initial callback
        // No need to thread
        f.listen();
    }
}

// Its a magic bus. Get it?
class MrsFrizzle {
    DataInputStream in;
    DataOutputStream out;
    Thread thread;

    MrsFrizzle(String host, int port) {
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
                    while (true) {
                        try {

                            // InputStream is = mSocket.getInputStream();
                            // DataInputStream dis = new DataInputStream(new GZIPInputStream(is));
                            int len = in.readInt();
                            byte[] buff = new byte[len];
                            in.readFully(buff);
                            String response = new String(buff, "UTF-8");
                            // Data data = new Gson().fromJson(response, Data.class);

//                            String string = "{\"id\":1,\"method\":\"object.deleteAll\",\"params\":[\"subscriber\"]}";
//
//                            PrintWriter pw = new PrintWriter(os);
//                            Main.log("Sent: " + string);
//                            pw.println(string);
//                            pw.flush();
//
//                            BufferedReader in = new BufferedReader(new InputStreamReader(is));
//                            String inputLine;
//
//                            while ((inputLine = in.readLine()) != null)
//                                System.out.println(inputLine);
//
//                            in.close();
//                            out.close();

                        } catch (Exception e) {
                            Main.log("Things didnt go well");
                            e.printStackTrace();
                        }
                    }
                }

        };

        thread.start();

    }

    public void sendResponse(String response) {
//        OutputStream os = mClientSocket.getOutputStream();
//        DataOutputStream dos = new DataOutputStream(new GZIPOutputStream(os));

        try {
            byte[] buff = response.getBytes("UTF-8");
            out.writeInt(buff.length);
            out.write(buff);
            out.flush();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}