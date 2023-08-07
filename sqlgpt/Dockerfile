# Use an official Python runtime as a parent image
FROM python:3.9

# Set the working directory in the container
WORKDIR /app

# Install any needed packages specified in requirements.txt
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Copy training.py to the container at /app
COPY main.py ./
RUN ls -halt

# Run the script when the container launches
CMD ["python", "main.py"]
