import os
from PIL import Image

def make_white_transparent(image_path, output_path):
    try:
        img = Image.open(image_path)
        img = img.convert("RGBA")
        datas = img.getdata()

        newData = []
        for item in datas:
            # Check if pixel is white or close to white
            if item[0] > 240 and item[1] > 240 and item[2] > 240:
                newData.append((255, 255, 255, 0)) # Transparent
            else:
                newData.append(item)

        img.putdata(newData)
        img.save(output_path, "PNG")
        print(f"Converted {image_path} to {output_path}")
        return True
    except Exception as e:
        print(f"Error processing {image_path}: {e}")
        return False

def process_directory(directory):
    for filename in os.listdir(directory):
        if filename.endswith(".jpg") or filename.endswith(".jpeg") or filename.endswith(".png"):
            filepath = os.path.join(directory, filename)
            # We will convert all to PNG. Even PNGs might have white backgrounds.
            name, _ = os.path.splitext(filename)
            output_filepath = os.path.join(directory, f"{name}.png")
            
            # If it's a JPG, remove the original JPG after conversion
            is_jpg = filename.endswith(".jpg") or filename.endswith(".jpeg")
            
            if make_white_transparent(filepath, output_filepath):
                if is_jpg and os.path.exists(filepath):
                    os.remove(filepath)
                    print(f"Removed original {filepath}")

if __name__ == "__main__":
    process_directory("web/public/components")
