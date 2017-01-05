create_cert(){
    openssl req \
        -new \
        -newkey rsa:4096 \
        -days 365 \
        -nodes \
        -x509 \
        -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com" \
        -keyout www.example.com.key \
        -out www.example.com.cert

    aws iam upload-server-certificate \
        --server-certificate-name l0-test-cert \
        --certificate-body www.example.com.cert \
        --private-key www.example.com.key
}

delete_cert() {
    rm www.example.com.key && rm www.example.com.cert

    aws iam delete-server-certificate
}

