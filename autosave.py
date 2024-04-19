from typing import List
import json
import os
from kitty.boss import Boss

def main(args: List[str]) -> str:
    pass

from kittens.tui.handler import result_handler
@result_handler(no_ui=True)
def handle_result(args: List[str], answer: str, target_window_id: int, boss: Boss) -> None:
    w = boss.window_id_map.get(target_window_id)
    if w is not None:
        json_data_str = boss.call_remote_control(w, ('ls', f'--match=id:{w.id}'))
        json_data = json.loads(json_data_str)
        # Define the file path
        home_dir = os.path.expanduser("~")
        
        # Define the file path with the correct username
        file_path = os.path.join(home_dir, '.config','kitty', 'lastsession.kitty')        
        # Check if the file exists
        if not os.path.exists(file_path):
            # If the file doesn't exist, create a new one
            with open(file_path, 'w') as f:
                json.dump({}, f)  # Write an empty JSON object
        
        # Write JSON data to the file
        with open(file_path, 'w') as f:
            json.dump(json_data, f)
