# todos/views.py
from django.http import JsonResponse, HttpResponseNotFound, HttpResponseServerError
import os
import json

def todos_view(request):
    try:
        db_path = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), '..', 'db.json')
        with open(db_path, 'r') as file:
            data = json.load(file)
        
        return JsonResponse(data)

    except FileNotFoundError:
        return HttpResponseNotFound('<h1>db.json file not found</h1>')

    except json.JSONDecodeError:
        return HttpResponseServerError('<h1>Invalid JSON format in db.json</h1>')

    except Exception as e:
        return HttpResponseServerError(f'<h1>Server Error: {str(e)}</h1>')
