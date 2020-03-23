until $(curl --silent --fail http://localhost:2216/health); do
    printf '.'
    sleep 5
done

echo "Starting cloudworker"
exec ./cloudworker