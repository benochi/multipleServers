from django.urls import path
from .views import todos_view

urlpatterns = [
    path('todos', todos_view, name='todos'),
]
