# 补录数据，csv转换成json
from datetime import datetime
import json, os
import pandas as pd
import sys

sys.path.append("./")
LOCAL_BASEDIR = ""


def is_valid_date_format(date_string, date_format="%Y-%m-%d"):
    try:
        datetime.strptime(date_string, date_format)
        return True
    except ValueError:
        print(ValueError)
        return False


def replace_nones(obj):
    if isinstance(obj, str):
        # obj = obj.replace("*", "")
        if is_valid_date_format(obj):
            return obj
        else:
            try:
                b = eval(
                    obj.replace(": None", ': ""')
                    .replace(": True", ": true")
                    .replace(": False", ": false")
                    .replace("*", "\*")
                )
                obj = json.loads(json.dumps(b))
            except Exception as e:
                print(e)
    if isinstance(obj, list):
        for i, elem in enumerate(obj):
            obj[i] = replace_nones(elem)
    elif isinstance(obj, dict):
        for key, value in obj.items():
            obj[key] = replace_nones(value)
    elif obj is None:
        return "null"
    return obj


LOCAL_BASEDIR = "/Users//temp/data"


def csv2json(folder_path):
    """
    Convert CSV files to JSON.
    Args:
        folder_path (str):
        The path of the folder containing CSV files.
    """
    folder_path = folder_path.replace(LOCAL_BASEDIR, "")
    # Get the paths of all CSV files in the folder
    local_folder_path = f"{LOCAL_BASEDIR}{folder_path}"
    csv_paths = [
        os.path.join(local_folder_path, f)
        for f in os.listdir(local_folder_path)
        if f.endswith(".csv")
    ]
    # Read all CSV files and merge them into one DataFrame
    df = pd.concat([pd.read_csv(f, dtype=str) for f in csv_paths], ignore_index=True)
    fill_df = df.fillna(value="")
    if "file_path" in fill_df.columns:
        fill_df.drop(columns="file_path", inplace=True)
    json_file_name = csv_paths[-1].replace(".csv", ".json")
    json_file_names = json_file_name.split("/")
    base_json_dir = "/".join(json_file_names[:-2])
    last_json_dir = "/".join(json_file_names[-2:])
    full_json_file_name = f"{base_json_dir}/json/{last_json_dir}"
    event_date_dir = f"{base_json_dir}/json/{json_file_names[-2]}"
    # Create the JSON directory if it doesn't exist
    if not os.path.exists(event_date_dir):
        os.makedirs(event_date_dir)
    # Iterate over each row of the DataFrame
    with open(full_json_file_name, "a") as f:
        for _, row in fill_df.iterrows():
            print(row)
            # Convert the row into JSON format
            row_as_json = json.dumps(replace_nones(row.to_dict()))
            # row.to_json(default_handler=replace_nones)
            # Open the file in append mode and write the data
            f.write(row_as_json)
            f.write("\n")  # Write a newline character to the file.
    return event_date_dir


base_dir = "/Users/.../every_30days/"

if __name__ == "__main__":
    for dirpath, dirnames, filenames in os.walk(base_dir):
        # # 打印当前目录下的所有文件夹
        for dirname in dirnames:
            csv2json(dirpath + dirname)
